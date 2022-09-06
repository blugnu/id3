package id3v23

import (
	"io"

	"github.com/blugnu/tags/internal/reader"
)

type ExtendedHeader struct {
	Size    uint32
	Padding uint32
	CRC     []byte
}

func ReadExtendedHeader(src io.ReadSeeker) (*ExtendedHeader, error) {
	eh := &extendedHeader{}
	err := eh.read(src)
	if err != nil {
		return nil, err
	}

	header := &ExtendedHeader{}
	header.Size = eh.size
	header.Padding = eh.padding
	header.CRC = eh.crc

	return header, nil
}

type extendedHeader struct {
	size    uint32 // [sync safe] size of the extended header (excluding itself)
	flags   []byte // extended header flags (2 bytes).  Currently only 1 flag is used: %x000 0000 0000 0000 - when x is set, a CRC-32 is also included in the header
	padding uint32 // [sync safe] size in bytes of any padding applied to the tag
	crc     []byte // CRC-32 (4 bytes) of frame data (before unsynchronisation), if present (see flags)
}

func (h *extendedHeader) read(src io.ReadSeeker) error {
	var err error

	reader := reader.New(src)

	h.size, err = reader.ReadSyncSafeUint32()
	if err != nil {
		return err
	}

	h.flags, err = reader.ReadBytes(2)
	if err != nil {
		return err
	}

	h.padding, err = reader.ReadSyncSafeUint32()
	if err != nil {
		return err
	}

	hasCrc := false
	if h.getFlags(&hasCrc); hasCrc {
		h.crc, err = reader.ReadBytes(4)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *extendedHeader) getFlags(crc *bool) {
	*crc = h.flags[0]&0x80 > 0
}
