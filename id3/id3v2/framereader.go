package id3v2

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v2/frame"
	"github.com/blugnu/tags/internal/reader"
)

type framereader struct {
	reader.Reader
	*Tag
	*frame.Frame
}

func (reader *framereader) readFrame() error {
	var id string
	var size int
	var flags uint16
	loc := reader.Pos()

	err := reader.readHeader(&id, &size, &flags)
	if err != nil {
		return err
	}

	reader.Frame = &frame.Frame{
		ID:       id,
		Size:     size,
		Location: loc,
	}
	parseFlags[reader.Tag.Version](reader.Frame, flags)

	switch reader.ID {

	case "PIC":
	case "APIC":
		err = reader.readPictureFrame()

	case "COM":
	case "COMM":
		err = reader.readCommentFrame()

	case "TXX":
	case "TXXX":
		err = reader.readUserDefinedTextFrame()

	default:
		if strings.HasPrefix(reader.ID, "T") {
			err = reader.readTextFrame()
		} else {
			err = reader.readUnknownFrame()
		}
	}
	if err != nil {
		reader.Frame = nil
		return err
	}

	return nil
}

func (reader *framereader) readCommentFrame() error {
	if err := reader.readTextEncoding(); err != nil {
		return err
	}
	if err := reader.readLanguageCode(); err != nil {
		return err
	}
	buf, err := reader.ReadBytes(reader.Frame.Size - 4) // TextEncoding + Language code = 4 bytes
	if err != nil {
		return err
	}

	el := bytes.Split(buf, reader.TextEncoding.Terminator())

	desc, err := reader.DecodeString(el[0])
	if err != nil {
		return err
	}
	comment, err := reader.DecodeString(el[1])
	if err != nil {
		return err
	}

	reader.Frame.Description = &desc
	reader.Frame.Text = &comment
	return nil
}

func (reader *framereader) readPictureFrame() error {
	if err := reader.readTextEncoding(); err != nil {
		return err
	}

	enc := frame.Iso88591
	buf, err := reader.ReadBytez(enc.Terminator())
	if err != nil {
		return err
	}
	mime, err := enc.Decode(buf)
	if err != nil {
		return err
	}
	lenImageType := len(buf) + 1

	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b > byte(frame.MaxPictureType) {
		return fmt.Errorf("unsupported picture type (%v)", b)
	}
	pictureType := frame.PictureType(b)

	enc = *reader.TextEncoding
	buf, err = reader.ReadBytez(enc.Terminator())
	if err != nil {
		return err
	}
	description, err := enc.Decode(buf)
	if err != nil {
		return err
	}
	lenDescription := len(buf) + len(enc.Terminator())

	data, err := reader.ReadBytes(reader.Frame.Size - (1 + lenImageType + 1 + lenDescription))
	if err != nil {
		return err
	}

	mime = strings.ToLower(mime)
	switch mime {
	case "gif":
		mime = "image/gif"
	case "png":
		mime = "image/png"
	case "jpg":
	case "jpeg":
		mime = "image/jpeg"
	}

	reader.Frame.Picture = &frame.Picture{
		MimeType:    mime,
		PictureType: pictureType,
		Description: description,
		Data:        data,
	}

	return nil
}

func (reader *framereader) readTextFrame() error {
	if err := reader.readTextEncoding(); err != nil {
		return err
	}

	buf, err := reader.ReadBytes(reader.Frame.Size - 1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}

	s, err := reader.DecodeString(buf)
	if err != nil {
		return err
	}
	reader.Text = &s
	return nil
}

func (reader *framereader) readUnknownFrame() error {
	data, err := reader.ReadBytes(reader.Frame.Size)
	if err != nil {
		return err
	}
	reader.Frame.UnknownData = data
	return nil

}

func (reader *framereader) readUserDefinedTextFrame() error {
	if err := reader.readTextEncoding(); err != nil {
		return err
	}
	buf, err := reader.ReadBytes(reader.Frame.Size - 1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}

	el := bytes.Split(buf, reader.TextEncoding.Terminator())

	desc, err := reader.DecodeString(el[0])
	if err != nil {
		return err
	}
	comment, err := reader.DecodeString(el[1])
	if err != nil {
		return err
	}

	reader.Frame.Description = &desc
	reader.Frame.Text = &comment
	return nil
}

func (reader *framereader) readHeader(id *string, size *int, flags *uint16) error {
	var err error

	// check for the presence of a valid frame id of the appropriate size
	var idlen = map[id3.TagVersion]int{
		id3.Id3v22: 3,
		id3.Id3v23: 4,
		id3.Id3v24: 5,
	}
	idbytes, err := reader.ReadBytes(idlen[reader.Tag.Version])
	if err != nil {
		return err
	}
	if !frame.IsValidId(idbytes) {
		return io.EOF
	}
	*id = string(idbytes)

	// we got a valid frame id, so now we can read the frame size...
	var sizefunc = map[id3.TagVersion]func() (uint32, error){
		id3.Id3v22: reader.ReadUint24,
		id3.Id3v23: reader.ReadUint32,
		id3.Id3v24: reader.ReadSyncSafeUint32,
	}
	ui32, err := sizefunc[reader.Tag.Version]()
	if err != nil {
		return err
	}
	*size = int(ui32)

	// ... for a 2.2.0 frame, that's it, but for 2.3.0 and later ...
	if reader.Tag.Version == id3.Id3v22 {
		return nil
	}

	// ... the header also contains 2 additional bytes of frame flags
	*flags, err = reader.ReadUint16()
	if err != nil {
		return err
	}

	return nil
}

func (reader *framereader) readLanguageCode() error {
	buf, err := reader.ReadBytes(3)
	if err != nil {
		return err
	}
	sbuf := string(buf)
	reader.Frame.LanguageCode = &sbuf
	return nil
}

func (reader *framereader) readTextEncoding() error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	reader.Frame.TextEncoding = frame.TextEncodingFromByte(b)
	return nil
}

var parseFlags = map[id3.TagVersion]func(frame *frame.Frame, flags uint16){
	id3.Id3v22: func(frame *frame.Frame, flags uint16) { /* NO-OP */ },
	id3.Id3v23: func(frame *frame.Frame, flags uint16) {
		// https://id3.org/id3v2.3.0#Frame_header_flags
		// %abc00000 ijk00000
		const tagAlterPreservationBit = 0x8000  // a
		const fileAlterPreservationBit = 0x4000 // b
		const readonlyBit = 0x2000              // c
		const compressionBit = 0x0080           // i
		const encryptionBit = 0x0040            // j
		const groupingBit = 0x0020              // k

		frame.PreserveWhenTagAltered = flags&tagAlterPreservationBit > 0
		frame.PreserveWhenFileAltered = flags&fileAlterPreservationBit > 0
		frame.IsReadOnly = flags&readonlyBit > 0
		frame.IsCompressed = flags&compressionBit > 0
		frame.IsEncrypted = flags&encryptionBit > 0
		frame.IsGrouped = flags&groupingBit > 0
	},
	id3.Id3v24: func(frame *frame.Frame, flags uint16) {
		// https://id3.org/id3v2.4.0-structure
		// %0abc0000 0h00kmnp
		const tagAlterPreservationBit = 0x4000  // a
		const fileAlterPreservationBit = 0x2000 // b
		const readonlyBit = 0x1000              // c
		const groupingBit = 0x0040              // h
		const compressionBit = 0x0008           // k
		const encryptionBit = 0x0004            // m
		const unsynchronisationBit = 0x0002     // n
		const datalengthBit = 0x0001            // p

		frame.PreserveWhenTagAltered = flags&tagAlterPreservationBit > 0
		frame.PreserveWhenFileAltered = flags&fileAlterPreservationBit > 0
		frame.IsReadOnly = flags&readonlyBit > 0
		frame.IsGrouped = flags&groupingBit > 0
		frame.IsCompressed = flags&compressionBit > 0
		frame.IsEncrypted = flags&encryptionBit > 0
		frame.IsUnsynchronised = flags&unsynchronisationBit > 0
		frame.HasDataLength = flags&datalengthBit > 0
	}}
