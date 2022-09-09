package v1filer

import (
	"strings"

	"github.com/blugnu/tags/internal/id3storage"
)

func (tag *reader) readV1String(n int) (string, error) {
	bytes := make([]byte, n)
	nr, err := tag.Read(bytes)
	if err != nil {
		return "", err
	}
	if nr < n {
		return "", id3storage.InsufficientData{Needed: n, Have: nr}
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
