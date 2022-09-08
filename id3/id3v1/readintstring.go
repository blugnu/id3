package id3v1

import (
	"strconv"
)

// Reads the specified number of bytes as a string and returns the
// result of parsing that string as an int.
func (reader *reader) readIntString(n int) (int, error) {
	s, err := reader.readString(4)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseUint(s, 10, 32)
	return int(v), err
}
