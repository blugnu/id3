package id3v2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/internal/reader"
)

var id3v2HeaderSIG = []byte("ID3")

// var id3v2FooterSIG = []byte("3DI")

type header struct {
	sig     []byte
	version struct {
		major    byte
		revision byte
	}
	flags   byte
	tagSize uint32
}

func (h *header) read(src io.ReadSeeker) error {
	var err error

	reader := reader.New(src)
	opos, _ := src.Seek(0, io.SeekCurrent)

	h.sig, err = reader.ReadBytes(3)
	if err != nil {
		return err
	}
	if !h.isValidHeader() {
		return id3.NoTag{AtPos: opos}
	}

	h.version.major, err = reader.ReadByte()
	if err != nil {
		return err
	}
	h.version.revision, err = reader.ReadByte()
	if err != nil {
		return err
	}
	h.flags, err = reader.ReadByte()
	if err != nil {
		return err
	}
	h.tagSize, err = reader.ReadSyncSafeUint32()
	if err != nil {
		return err
	}

	// In the unlikely event that we have a v2.2.0 tag with the compression flag set
	// we return an UnsupportedTag error
	if h.version.major == 2 && h.flags&0x40 > 0 {
		return id3.UnsupportedTag{Reason: "2.2.0 tag with compression"}
	}

	// In the similarly unlikely even of encountering an IDv3 tag < v2.2.0 or > v2.4.0
	// that also is an UnsupportedTag
	if h.version.major < 2 || h.version.major > 4 {
		return id3.UnsupportedTag{Reason: fmt.Sprintf("ID3v2.%d.%d tag not supported", h.version.major, h.version.revision)}
	}

	return nil
}

func (h *header) getFlags(a *bool, b *bool, c *bool, d *bool) {
	*a = h.flags&0x80 > 0
	*b = h.flags&0x40 > 0
	*c = h.flags&0x20 > 0
	*d = h.flags&0x10 > 0
}

func (h *header) getVersion() id3.TagVersion {
	switch h.version.major {
	case 0x02:
		return id3.Id3v22
	case 0x03:
		return id3.Id3v23
	case 0x04:
		return id3.Id3v24
	default:
		return id3.Id3vUnknown
	}
}

func (h *header) isValidHeader() bool {
	return bytes.Equal(h.sig, id3v2HeaderSIG)
}

// func (h *header) size() uint { return uint(unsafe.Sizeof(header{})) }
