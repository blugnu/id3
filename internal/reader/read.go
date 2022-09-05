package reader

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type reader struct {
	io.Reader
}

func New(src io.Reader) *reader {
	return &reader{src}
}

// Reads and returns a single byte
func (r *reader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r.Reader, b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// Reads the specified number of bytes and returns the result in a slice
func (r *reader) ReadBytes(n uint) ([]byte, error) {
	const max = 10 << 20

	if n > max {
		b := &bytes.Buffer{}
		if _, err := io.CopyN(b, r.Reader, int64(n)); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}

	b := make([]byte, n)
	_, err := io.ReadFull(r.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Reads the specified number of bytes and returns the result as a
// string (any trailing 0x00 and leading/trailing whitespace removed)
func (r *reader) ReadString(n uint) (string, error) {
	bytes, err := r.ReadBytes(n)
	if err != nil {
		return "", err
	}

	// Remove any zero padding
	for bytes[len(bytes)-1] == 0 {
		bytes = bytes[:len(bytes)-1]
	}

	return strings.TrimSpace(string(bytes)), nil
}

// Reads the specified number of bytes as a string and returns the
// result of parsing that string as an int.
func (r *reader) ReadStringAsInt(n uint) (int, error) {
	b, err := r.ReadBytes(n)
	if err != nil {
		return 0, err
	}
	s := string(b)
	v, err := strconv.ParseUint(s, 10, 32)
	return int(v), err
}

// Reads and Decodes a 4-byte sync-safe Uint32 value
func (r *reader) ReadSyncSafeUint32() (uint32, error) {
	buf, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	// Sync-safe ints are stored in MSB order, so on a Little Endian
	// platform we need to reverse the bytes
	if byteOrder == binary.LittleEndian {
		buf = reverse(buf)
	}

	// Now can unsync the value
	return unsync(buf)
}
