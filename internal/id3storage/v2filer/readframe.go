package v2filer

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v2"
)

// signature type for funcs returning frame size
type readFrameSize func() (uint32, error)

type framereader struct {
	*reader       // ref to an id3 byte-level reader
	readFrameSize // ref to the function used to read frame size
	*id3v2.Tag    // ref to the tag for which the reader is reading frames
	*id3v2.Frame  // a ref to the single frame obtained by an individual call to reader.readFrame()
}

func (reader *framereader) readFrame() error {
	var id string
	var size int
	var flags uint16

	err := reader.readHeader(&id, &size, &flags)
	if err == io.EOF {
		return err
	}
	if err != nil {
		return fmt.Errorf("readFrame [header]: %w", err)
	}

	key := id3.FrameKeyForID(id)

	reader.Frame = &id3v2.Frame{
		ID:   id,
		Key:  key,
		Size: size,
	}

	switch reader.Tag.Version {
	case id3.Id3v23:
		reader.Frame.Flags = &id3v2.FrameFlags{
			PreserveWhenTagAltered:  flags&frameheader.v230flag.tagAlterPreservation > 0,
			PreserveWhenFileAltered: flags&frameheader.v230flag.fileAlterPreservation > 0,
			IsReadOnly:              flags&frameheader.v230flag.readonly > 0,
			IsCompressed:            flags&frameheader.v230flag.compression > 0,
			IsEncrypted:             flags&frameheader.v230flag.encryption > 0,
			IsGrouped:               flags&frameheader.v230flag.grouping > 0,
		}
	case id3.Id3v24:
		reader.Frame.Flags = &id3v2.FrameFlags{
			PreserveWhenTagAltered:  flags&frameheader.v240flag.tagAlterPreservation > 0,
			PreserveWhenFileAltered: flags&frameheader.v240flag.fileAlterPreservation > 0,
			IsReadOnly:              flags&frameheader.v240flag.readonly > 0,
			IsGrouped:               flags&frameheader.v240flag.grouping > 0,
			IsCompressed:            flags&frameheader.v240flag.compression > 0,
			IsEncrypted:             flags&frameheader.v240flag.encryption > 0,
			IsUnsynchronised:        flags&frameheader.v240flag.unsynchronisation > 0,
			HasDataLength:           flags&frameheader.v240flag.datalength > 0,
		}
	}

	switch reader.Frame.Key {

	case id3.UnknownKey:
		reader.Frame = nil
		return fmt.Errorf("readFrame [%s]: unknown id", id)

	case id3.APIC:
		err = reader.readPictureFrame()

	case id3.COMM:
		err = reader.readCommentFrame()

	case id3.TPOS, id3.TRCK:
		err = reader.readPartOfSetFrame()

	case id3.TXXX:
		err = reader.readUserDefinedTextFrame()

	default:
		if id[0] == 'T' {
			err = reader.readTextFrame()
		} else {
			err = reader.readUnknownFrame()
		}
	}
	if err != nil {
		reader.Frame = nil
		return fmt.Errorf("readFrame [%s]: %w", id, err)
	}

	return nil
}

func (reader *framereader) readHeader(id *string, size *int, flags *uint16) error {
	var err error

	idlen := frameidlen[reader.Tag.Version]
	idbytes, _, err := reader.readBytes(idlen)
	if err != nil {
		return fmt.Errorf("readHeader [id]: %w", err)
	}
	if isNullFrameId(idbytes) {
		return io.EOF
	}
	if !isValidFrameId(idbytes) {
		return fmt.Errorf("not a valid frame id (%s)", idbytes)
	}
	*id = string(idbytes)

	// we got a valid frame id, so now we can read the frame size...
	ui32, err := reader.readFrameSize()
	if err != nil {
		return fmt.Errorf("readHeader [frameSize]: %w", err)
	}
	*size = int(ui32)

	// ... for a 2.2.0 frame the header is complete...
	if reader.Tag.Version == id3.Id3v22 {
		return nil
	}

	// ... but for 2.30 and 2.40 the header also contains 2 additional
	// bytes of frame flags.  We'll need to parse those later, once we
	// have a frame to parse them into; for now we just set the flag
	// bytes before returning
	*flags, err = reader.readUint16()
	if err != nil {
		return fmt.Errorf("readHeader [flags]: %w", err)
	}

	return nil
}

func (frame *framereader) readCommentFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}
	lang, err := frame.readLanguageCode()
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	d, dlen, err := frame.readStringz(enc)
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}
	c, _, err := frame.readString(enc, frame.Frame.Size-4-dlen) // TextEncoding + Language code = 4 bytes
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	frame.Data = &id3v2.Comment{
		LanguageCode: lang,
		Description:  d,
		Comment:      c,
	}

	return nil
}

func (frame *framereader) readPictureFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	mime, nmimebytes, err := frame.readStringz(Iso88591) //
	if err != nil {
		return err
	}

	b, err := frame.readByte()
	if err != nil {
		return err
	}
	if b > byte(id3v2.StudioLogo) {
		return fmt.Errorf("unsupported picture type (%v)", b)
	}
	pictureType := id3v2.PictureType(b)

	description, ndescbytes, err := frame.readStringz(enc)
	if err != nil {
		return err
	}

	data, _, err := frame.readBytes(frame.Frame.Size - (nmimebytes + ndescbytes + 2)) // +2 = 1 each for text encoding and picture type
	if err != nil {
		return err
	}

	mime = strings.ToLower(mime)
	switch mime {
	case "gif":
		mime = "image/gif"
	case "png":
		mime = "image/png"
	case "jpg", "jpeg":
		mime = "image/jpeg"
	}

	frame.Data = &id3v2.Picture{
		MimeType:    mime,
		PictureType: pictureType,
		Description: description,
		Data:        data,
	}

	return nil
}

func (reader *framereader) readLanguageCode() (string, error) {
	buf, _, err := reader.readBytes(3)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (reader *framereader) readPartOfSetFrame() error {
	err := reader.readTextFrame()
	if err != nil {
		return fmt.Errorf("readPartOfSetFrame: %w", err)
	}
	s := reader.Frame.Data.(string)
	el := strings.Split(s, "/")
	var n int64 = -1
	var c int64 = -1
	switch len(el) {
	case 0:
		// NO-OP (will be set as '-1 of -1')

	case 1:
		n, err = strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("readPartOfSetFrame [item no]: %w", err)
		}
		c = -1

	case 2:
		if len(el[0]) > 0 {
			n, err = strconv.ParseInt(el[0], 10, 32)
			if err != nil {
				return fmt.Errorf("readPartOfSetFrame [item no]: %w", err)
			}
		}
		if len(el[1]) > 0 {
			c, err = strconv.ParseInt(el[1], 10, 32)
			if err != nil {
				return fmt.Errorf("readPartOfSetFrame [item count]: %w", err)
			}
		}

	default:
		return fmt.Errorf("readPartOfSetFrame: invalid (%v)", s)
	}

	reader.Frame.Data = &id3v2.PartOfSet{
		ItemNo:    int(n),
		ItemCount: int(c),
	}

	return nil
}

func (reader *framereader) readTextEncoding() (TextEncoding, error) {
	const invalid = 0xff

	b, err := reader.readByte()
	if err != nil {
		return invalid, err
	}

	if b > byte(Utf8) {
		return invalid, fmt.Errorf("invalid TextEncoding (%x)", b)
	}
	return TextEncoding(b), nil
}

func (reader *framereader) readString(enc TextEncoding, maxlen int) (string, int, error) {
	// read the maxlen bytes specified
	buf, n, err := reader.readBytes(maxlen)
	if err != nil {
		return "", n, fmt.Errorf("ReadString: %w", err)
	}

	// remove any null-terminator(s) for the specified encoding
	zlen := enc.zlen()
	zterm := make([]byte, zlen)
	for len(buf) > 0 && bytes.Equal(buf[len(buf)-zlen:], zterm) {
		buf = buf[:len(buf)-(zlen+1)]
	}

	// decode the buffer using the specified encoding
	s, err := enc.decode(buf)
	if err != nil {
		return "", n, fmt.Errorf("ReadString: %w", err)
	}

	// return the decoded string and the number of bytes originally read (including any null terminators)
	return s, n, nil
}

func (reader *framereader) readStringz(enc TextEncoding) (string, int, error) {

	zlen := enc.zlen()

	var buf []byte
	var err error
	if zlen == 1 {
		buf, err = reader.readBytez()
	} else {
		buf, err = reader.readBytezz()
	}
	if err != nil {
		return "", 0, fmt.Errorf("ReadStringz: %w", err)
	}
	readlen := len(buf) + zlen

	s, err := enc.decode(buf)
	if err != nil {
		return "", readlen, fmt.Errorf("ReadStringz: %w", err)
	}

	return s, readlen, nil
}

func (frame *framereader) readTextFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	// TODO: multiple null-termed values...

	str, _, err := frame.readString(enc, frame.Frame.Size-1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}
	frame.Data = str

	return nil
}

func (frame *framereader) readUnknownFrame() error {
	data, _, err := frame.readBytes(frame.Frame.Size)
	if err != nil {
		return err
	}

	frame.Data = data

	return nil

}

func (frame *framereader) readUserDefinedTextFrame() error {

	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	d, dlen, err := frame.readStringz(enc)
	if err != nil {
		return fmt.Errorf("readUserDefinedText: %w", err)
	}

	t, _, err := frame.readString(enc, frame.Frame.Size-1-dlen)
	if err != nil {
		return fmt.Errorf("readUserDefinedText: %w", err)
	}

	frame.Data = &id3v2.UserDefinedText{
		Description: d,
		Text:        t,
	}

	return nil
}
