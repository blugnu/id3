package id3v24

import (
	"fmt"
	"io"

	"github.com/blugnu/tags/internal/reader"
)

type TagSizeRestriction byte
type TextEncodingRestriction byte
type ImageEncodingRestriction byte
type ImageSizeRestriction byte

type ExtendedHeader struct {
	Size         uint32
	IsUpdate     bool
	CRC          *uint32
	Restrictions struct {
		TagSize       TagSizeRestriction
		TextEncoding  TextEncodingRestriction
		ImageEncoding ImageEncodingRestriction
		ImageSize     ImageSizeRestriction
	}
}

func ReadExtendedHeader(src io.ReadSeeker) (*ExtendedHeader, error) {
	eh := &extendedHeader{}
	err := eh.read(src)
	if err != nil {
		return nil, err
	}

	header := &ExtendedHeader{}
	header.Size = eh.size

	return header, nil
}

type extendedHeader struct {
	size     uint32 // [sync safe] size of the extended header (excluding itself)
	isUpdate bool
	crc      uint32
}

func (h *extendedHeader) read(src io.ReadSeeker) error {
	var err error

	reader := reader.New(src)

	h.size, err = reader.ReadSyncSafeUint32()
	if err != nil {
		return err
	}

	flagBytes, err := reader.ReadByte()
	if err != nil {
		return err
	}

	epos, _ := src.Seek(0, io.SeekCurrent)
	epos += int64(flagBytes)

	for {
		err := h.readFlag(reader)
		if err != nil {
			src.Seek(epos, io.SeekStart)
			return err
		}
		p, _ := src.Seek(0, io.SeekCurrent)
		if p == epos {
			break
		}
	}

	return nil
}

// Reads a flag and returns the number of bytes of data used by that flag
func (h *extendedHeader) readFlag(reader reader.Reader) error {
	const updateBit = 0x40
	const crcBit = 0x20
	const restrictionsBit = 0x10

	// Every flag consists of:
	// - a byte containing the flag-type indicator
	// - a byte identifying the length of any additional data for the flag
	// - any additional data bytes (or none)

	// Read the flag type indicator
	flag, err := reader.ReadByte()
	if err != nil {
		return err
	}

	// update flag
	if flag&updateBit > 0 {
		h.isUpdate = true
		_, err := h.readFlagData(reader, 0, "update")
		if err != nil {
			return err
		}
		return nil
	}

	// crc flag
	if flag&crcBit > 0 {
		crc, err := h.readFlagData(reader, 5, "crc")
		if err != nil {
			return err
		}
		h.crc, err = reader.UnsyncUint32(crc)
		if err != nil {
			return err
		}
		return nil
	}

	// restrictions flag
	if flag&restrictionsBit > 0 {
		_, err := h.readFlagData(reader, 1, "restrictions")
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("unsupported flag type: %x", flag)
}

// Reads the data bytes for a flag and returns an error if the expected number of data bytes does not match the
// data bytes indicated by the flag
func (h *extendedHeader) readFlagData(reader reader.Reader, expectedDataBytes byte, name string) ([]byte, error) {
	dataBytes, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if dataBytes != expectedDataBytes {
		return nil, fmt.Errorf("unexpected length for %s flag data (got %d, expected %d)", name, dataBytes, expectedDataBytes)
	}
	if dataBytes > 0 {
		return reader.ReadBytes(uint(dataBytes))
	}
	return nil, nil
}
