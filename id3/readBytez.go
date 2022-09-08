package id3

import (
	"fmt"
)

// Reads bytes until the specified terminator is encountered.
// Note that this is NOT reading a terminated string only a
// terminated byte buffer.
func (r *reader) ReadBytez(zn int) ([]byte, error) {
	switch zn {
	case 1:
		return r.readBytez()
	case 2:
		return r.readBytezz()
	default:
		return nil, fmt.Errorf("invalid null terminator length (%d): must be 1 or 2", zn)
	}
}

// Reads bytes until a null (Z-ero) is encountered.
// Note that this is NOT reading a null-terminated string only a
// null-terminated byte buffer.
func (r *reader) readBytez() ([]byte, error) {
	buf := []byte{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("readBytez: %w", err)
		}
		if b == 0x00 {
			break
		}
		buf = append(buf, b)
	}
	return buf, nil
}

// Reads bytes until a double null (Z-ero) is encountered.
// Note that this is NOT reading a null-terminated string only a
// double null-terminated byte buffer.
func (r *reader) readBytezz() ([]byte, error) {
	buf := []byte{}
	for {
		chunk, err := r.ReadBytes(2)
		if err != nil {
			return nil, fmt.Errorf("readBytezz: %w", err)
		}
		if chunk[0] == 0 && chunk[1] == 0 {
			break
		}
		buf = append(buf, chunk...)
	}
	return buf, nil
}
