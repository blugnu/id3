package tags

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

func readByte(src io.Reader) (byte, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(src, b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func readBytes(src io.Reader, n uint) ([]byte, error) {
	const max = 10 << 20

	if n > max {
		b := &bytes.Buffer{}
		if _, err := io.CopyN(b, src, int64(n)); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}

	b := make([]byte, n)
	_, err := io.ReadFull(src, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readString(r io.Reader, n uint) (string, error) {
	bytes, err := readBytes(r, n)
	if err != nil {
		return "", err
	}

	// Remove any zero padding
	for bytes[len(bytes)-1] == 0 {
		bytes = bytes[:len(bytes)-1]
	}

	return strings.TrimSpace(string(bytes)), nil
}

func readStringAsInt(r io.Reader, n uint) (int, error) {
	b, err := readBytes(r, n)
	if err != nil {
		return 0, err
	}
	s := string(b)
	v, err := strconv.ParseUint(s, 10, 32)
	return int(v), err
}
