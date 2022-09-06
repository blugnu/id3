package frame

import (
	"errors"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

func (frame *Frame) DecodeString(buf []byte) (string, error) {
	if frame.TextEncoding == nil {
		return "", errors.New("no text encoding")
	}

	switch *frame.TextEncoding {
	case Iso88591:
		dec := charmap.ISO8859_1.NewDecoder()
		decoded, err := dec.Bytes(buf)
		if err != nil {
			return "", err
		}
		return string(decoded), nil

	case Utf16:
		return "", errors.New("utf16 support not yet implemented")

	case Utf16BE:
		return "", errors.New("utf16BE support not yet implemented")

	case Utf8:
		return string(buf), nil
	default:
		return "", fmt.Errorf("text encoding (%d) not supported", *frame.TextEncoding)
	}
}
