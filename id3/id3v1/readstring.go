package id3v1

import (
	"strings"
)

func (tag *reader) readString(n int) (string, error) {
	bytes, err := tag.ReadBytes(n)
	if err != nil {
		return "", err
	}

	// Remove any zero padding
	switch len(bytes) {
	case 0:
		return "", nil

	case 1:
		if bytes[0] == 0 || bytes[0] == ' ' {
			return "", nil
		}
		return string(bytes), nil

	default:
		for len(bytes) > 0 && bytes[len(bytes)-1] == 0 {
			bytes = bytes[:len(bytes)-1]
		}
		return strings.TrimSpace(string(bytes)), nil
	}
}
