package id3v2

import (
	"bytes"
	"fmt"

	id3reader "github.com/blugnu/tags/internal/id3/reader"
)

func (header *tagreader) readHeader(pos *int64, majorver *byte, revision *byte, flags *byte, size *uint32) error {
	var err error

	*pos = header.Pos()
	sig, err := header.ReadBytes(3)
	if err != nil {
		return err
	}
	if !bytes.Equal(sig, []byte(id3v2HeaderSIG)) {
		return id3reader.NoTag{AtPos: *pos}
	}

	*majorver, err = header.ReadByte()
	if err != nil {
		return err
	}
	*revision, err = header.ReadByte()
	if err != nil {
		return err
	}
	*flags, err = header.ReadByte()
	if err != nil {
		return err
	}
	*size, err = header.ReadSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readHeader: %w", err)
	}

	// In the unlikely even of encountering an IDv3 tag < v2.2.0 or > v2.4.0
	// that is an UnsupportedTag
	if *majorver < 2 || *majorver > 4 {
		return id3reader.UnsupportedTag{Reason: fmt.Sprintf("ID3v2.%d.%d tag not supported", majorver, revision)}
	}

	// v2.2.0 tags with the compression are also unsupported
	if *majorver == 2 && *flags&headerflags.compression > 0 {
		return id3reader.UnsupportedTag{Reason: "2.2.0 tag with compression"}
	}

	return nil
}
