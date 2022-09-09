package v2filer

import (
	"bytes"
	"fmt"
	"io"
)

type reader struct {
	io.Reader
}

// Reads and returns a single byte
func (r *reader) readByte() (byte, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r.Reader, b)
	// n, err := io.ReadFull(r.Reader, b)
	// r.curr += int64(n)
	if err != nil {
		return 0, fmt.Errorf("ReadByte: %w", err)
	}
	return b[0], nil
}

// Reads the specified number of bytes and returns the result in a slice
func (r *reader) readBytes(n int) ([]byte, error) {
	if n == 0 {
		return []byte{}, nil
	}

	const max = 10 << 20

	if n > max {
		b := &bytes.Buffer{}
		if _, err := io.CopyN(b, r.Reader, int64(n)); err != nil {
			return nil, fmt.Errorf("ReadBytes (> max): %w", err)
		}
		// r.curr += int64(n)
		// if r.curr > r.max {
		// 	r.curr = r.max
		// }
		return b.Bytes(), nil
	}

	b := make([]byte, n)
	_, err := io.ReadFull(r.Reader, b)
	// nr, err := io.ReadFull(r.Reader, b)
	// r.curr += int64(nr)
	if err != nil {
		return nil, fmt.Errorf("ReadBytes: %w", err)
	}
	return b, nil
}

// Reads bytes until a null (Z-ero) is encountered.
// Note that this is NOT reading a null-terminated string only a
// null-terminated byte buffer.
func (r *reader) readBytez() ([]byte, error) {
	buf := []byte{}
	for {
		b, err := r.readByte()
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
		chunk, err := r.readBytes(2)
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

// Reads and Decodes a 4-byte sync-safe int
func (r *reader) readSyncSafeUint32() (uint32, error) {
	buf, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}

	// id3 stipulates that ints are stored in MSB order and the bit-shifting that is
	// performed when unsync-ing works as intended on a little-endian architecture.
	// this is UNTESTED on BIG-endian architectures, so if using this code in a
	// big-endian world, be prepared for trouble!

	v, err := unsyncUint32(buf)
	if err != nil {
		return 0, fmt.Errorf("ReadSyncSafeUint32: %w", err)
	}

	return v, nil
}

// Reads a uint16 from a 2-byte field
func (r *reader) readUint16() (uint16, error) {
	ui, err := r.readUint(2)
	if err != nil {
		return 0, err
	}
	return uint16(ui), nil
}

// Reads a uint32 from a 3-byte field
func (r *reader) readUint24() (uint32, error) {
	ui, err := r.readUint(3)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint32 from a 4-byte field
func (r *reader) readUint32() (uint32, error) {
	ui, err := r.readUint(4)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint64 from a field of n bytes (maximum 8)
func (r *reader) readUint(n int) (uint, error) {
	var result uint64

	bytes, err := r.readBytes(n)
	if err != nil {
		return 0, err
	}

	// id3 stipulates that ints are stored in MSB order and the bit-shifting below
	// works as intended on a little-endian architecture.  however, this is UNTESTED
	// on BIG-endian architectures, so if using this code in a big-endian world,
	// be prepared for trouble!

	for _, b := range bytes {
		result = (result << 8) | uint64(b)
	}
	return uint(result), nil
}

// seek provides convenient access to seek'ing the underlying io.Reader
// in the *reader
func (r *reader) seek(pos int64, from int) (int64, error) {
	return r.Reader.(io.ReadSeeker).Seek(pos, from)
}
