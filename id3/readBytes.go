package id3

import (
	"bytes"
	"fmt"
	"io"
)

// Reads and returns a single byte
func (r *reader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	n, err := io.ReadFull(r.Reader, b)
	r.curr += int64(n)
	if err != nil {
		return 0, fmt.Errorf("ReadByte: %w", err)
	}
	return b[0], nil
}

// Reads the specified number of bytes and returns the result in a slice
func (r *reader) ReadBytes(n int) ([]byte, error) {
	if n == 0 {
		return []byte{}, nil
	}

	const max = 10 << 20

	if n > max {
		b := &bytes.Buffer{}
		if _, err := io.CopyN(b, r.Reader, int64(n)); err != nil {
			return nil, fmt.Errorf("ReadBytes (> max): %w", err)
		}
		r.curr += int64(n)
		if r.curr > r.max {
			r.curr = r.max
		}
		return b.Bytes(), nil
	}

	b := make([]byte, n)
	nr, err := io.ReadFull(r.Reader, b)
	r.curr += int64(nr)
	if err != nil {
		return nil, fmt.Errorf("ReadBytes: %w", err)
	}
	return b, nil
}
