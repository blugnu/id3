package v2filer

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v2"
	"github.com/blugnu/tags/internal/id3storage"
)

type tagreader struct {
	*reader
	*id3v2.Tag
}

func ReadTag(src io.ReadSeeker) (*id3v2.Tag, error) {

	tr := &tagreader{&reader{src}, nil}

	if err := tr.readTag(); err != nil {
		if errors.Is(err, id3storage.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}

	return tr.Tag, nil
}

func (tag *tagreader) readTag() error {

	var pos int64
	var majorver byte
	var revision byte
	var flags byte
	var size uint32

	err := tag.readHeader(&pos, &majorver, &revision, &flags, &size)
	defer func() {
		// a non-zero size means we at least THINK we read a valid header
		// and should make sure we position at the end of the tag when
		// we're done, whatever happens
		if size > 0 {
			tag.seek(pos+tagHeaderSize+int64(size), io.SeekStart)
		}
	}()

	// if we get an UnsupportedTag error, this indicates we read a
	// header that was identified as a v2 tag ("ID3") but is of a version
	// or uses features that are not supported by this implementation
	//
	// in which case we identify it as a tag where the only known
	// information is the location, the size (excluding header) and
	// raw tag data (INCLUDING the entire header)
	if err != nil {
		if errors.Is(err, id3storage.UnsupportedTag{}) {
			tag.seek(pos, io.SeekStart)
			data, _, err := tag.readBytes(int(size) + tagHeaderSize)
			if err != nil {
				return err
			}
			tag.Tag = &id3v2.Tag{
				Location: pos,
				Size:     size,
				RawData:  data,
			}
		}
		return err
	}

	// otherwise we have a valid, supported tag...

	tag.Tag = &id3v2.Tag{
		Version:           tagVersion[majorver],
		Size:              size,
		Location:          pos,
		IsExperimental:    flags&headerflags.unsynchronisation > 0,
		IsUnsynchronised:  flags&headerflags.extendedHeader > 0,
		HasExtendedHeader: flags&headerflags.experimental > 0,
		HasFooter:         flags&headerflags.footer > 0,
	}

	// .. which may also have an extended header

	if err := tag.readExtendedHeader(); err != nil {
		return err
	}

	// .. and will definitely have at least one frame

	framedata, _, err := tag.readBytes(int(tag.Size))
	if err != nil {
		return err
	}

	// if tag.IsUnsynchronised {
	// TODO: apply de-unsynchronisation to tagdata before reading any frames
	// }

	err = tag.readFrames(framedata)
	if err != nil {
		return err
	}

	return nil
}

// readHeader reads an ID3v2 tag header at the current seek position.
func (header *tagreader) readHeader(pos *int64, majorver *byte, revision *byte, flags *byte, size *uint32) error {
	var err error

	*pos, _ = header.seek(0, io.SeekCurrent)
	sig, _, err := header.readBytes(3)
	if err != nil {
		return err
	}
	if !bytes.Equal(sig, []byte(tagHeaderSIG)) {
		return id3storage.NoTag{AtPos: *pos}
	}

	*majorver, err = header.readByte()
	if err != nil {
		return err
	}
	*revision, err = header.readByte()
	if err != nil {
		return err
	}
	*flags, err = header.readByte()
	if err != nil {
		return err
	}
	*size, err = header.readSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readHeader: %w", err)
	}

	// In the unlikely even of encountering an IDv3 tag < v2.2.0 or > v2.4.0
	// that is an UnsupportedTag
	if *majorver < 2 || *majorver > 4 {
		return id3storage.UnsupportedTag{Reason: fmt.Sprintf("ID3v2.%d.%d tag not supported", majorver, revision)}
	}

	// v2.2.0 tags with the compression are also unsupported
	if *majorver == 2 && *flags&headerflags.compression > 0 {
		return id3storage.UnsupportedTag{Reason: "2.2.0 tag with compression"}
	}

	return nil
}

func (tag *tagreader) readExtendedHeader() error {
	if !tag.HasExtendedHeader {
		return nil
	}

	switch tag.Version {
	case id3.Id3v23:
		return tag.readExtendedHeader230()
	case id3.Id3v24:
		return tag.readExtendedHeader240()
	}

	return nil
}

func (tag *tagreader) readExtendedHeader230() error {
	if _, err := tag.readSyncSafeUint32(); err != nil {
		return fmt.Errorf("readExtendedHeader230: %w", err)
	}

	flags, err := tag.readUint16()
	if err != nil {
		return err
	}

	tag.Padding, err = tag.readSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readExtendedHeader230: %w", err)
	}

	// if the CRC bit is set, read the CRC data.  In a 2.3.0 extended
	// header, this is a regular uint32, not sync-safe!
	if flags&extendedheader.v230flag.crc > 0 {
		tag.CRC, err = tag.readUint32()
		if err != nil {
			return err
		}
	}

	return nil
}

func (tag *tagreader) readExtendedHeader240() error {
	var err error

	pos, _ := tag.seek(0, io.SeekCurrent)
	size, err := tag.readSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readExtendedHeader240: %w", err)
	}

	// Ensure we finish with the reader positioned immediately after
	// the indicated extended header data
	defer func() {
		tag.seek(pos+int64(size), io.SeekStart)
	}()

	// Id3v2.4 future-proofs the extended header by allowing for a variable number
	// of flag bytes, but in 2.4.0 itself there is always only 1.
	//
	// We still need to read the flag length (and should check it to ensure it is the expected value)
	flagbytes, err := tag.readByte()
	if err != nil {
		return err
	}
	if flagbytes != 1 {
		return fmt.Errorf("found %d flag bytes (expected 1)", flagbytes)
	}

	// read the single flag byte
	flags, err := tag.readByte()
	if err != nil {
		return err
	}

	// interpret data for any set flags in order of the flag bits
	for _, bit := range []byte{
		extendedheader.v240flag.update,
		extendedheader.v240flag.crc,
		extendedheader.v240flag.restrictions,
	} {
		// if the required flag bit is not set then there is nothing to read
		if flags&bit == 0 {
			continue
		}
		err = tag.readFlag(bit)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reads data for the flag identified by the flag bit specified
func (tag *tagreader) readFlag(bit byte) error {
	// Every flag consists of:
	// - a byte identifying the length of any additional data for the flag
	// - any additional data bytes (or none)

	datalen, err := tag.readByte()
	if err != nil {
		return err
	}
	expected := v240flaginfo[bit]
	if datalen != expected.datalen {
		return fmt.Errorf(" %s flag: unexpected data length: got %d, expected %d", expected.name, datalen, expected.datalen)
	}
	data, _, err := tag.readBytes(int(datalen))
	if err != nil {
		return err
	}

	// now parse the flag data according to the flag type
	switch bit {
	// update flag - no data, it's mere presence IS the data
	case extendedheader.v240flag.update:
		tag.IsUpdate = true
		return nil

	// crc flag - has 5 bytes of data containing a 32-bit sync-safe int
	case extendedheader.v240flag.crc:
		tag.CRC, err = unsyncUint32(data)
		if err != nil {
			return err
		}
		return nil

		// restrictions flag - 1 bit of additional bits signifiying varying restrictions
	case extendedheader.v240flag.restrictions:
		// TODO: parse the bits and set restriction indicators in the tag (to be defined)
		// For now, having read the data byte is enough to "handle" this flag so we'll just
		// set the tag.Restrictions to the unparsed byte
		tag.Restrictions = data[0]
		return nil
	}

	return fmt.Errorf("flag bit not supported (%x)", bit)
}

// readFrames reads as many frames as can be found in the supplied framedata
// byte slice.  Frames are appended to the []Frames of the current *Tag in
// the *tagreader.
func (tag *tagreader) readFrames(framedata []byte) error {

	// careful.. this is a bit cute...
	//
	// when hooking up the func in the framereader to be used for ReadFrameSize,
	// we need to make sure we hookup the func provided by the framereader itself.
	//
	// so we initially get a framereader with NO ReadFrameSize func set ...
	reader := &framereader{
		reader:        &reader{bytes.NewReader(framedata)},
		readFrameSize: nil,
		Tag:           tag.Tag,
		Frame:         nil,
	}
	// .. and now we can give the framereader a ReadFrameSize func ref to one
	// of the funcs it provides itself!
	reader.readFrameSize = map[id3.TagVersion]readFrameSize{
		id3.Id3v22: reader.readUint24,
		id3.Id3v23: reader.readUint32,
		id3.Id3v24: reader.readSyncSafeUint32,
	}[tag.Version]

	for {
		err := reader.readFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		tag.Frames = append(tag.Frames, reader.Frame)
	}
	return nil
}
