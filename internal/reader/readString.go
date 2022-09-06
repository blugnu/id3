package reader

import (
	"strconv"
	"strings"
)

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
