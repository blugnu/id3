package v2filer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf16"

	"golang.org/x/text/encoding/charmap"
)

type TextEncoding byte

const (
	Iso88591 TextEncoding = iota
	Utf16
	Utf16BE
	Utf8
)

func (enc TextEncoding) decode(buf []byte) (string, error) {
	switch enc {
	case Iso88591:
		dec := charmap.ISO8859_1.NewDecoder()
		decoded, err := dec.Bytes(buf)
		if err != nil {
			return "", err
		}
		return string(decoded), nil

	case Utf16:
		wide, err := toUint16s(buf)
		if err != nil {
			return "", err
		}
		decoded := utf16.Decode(wide)
		return string(decoded), nil

	case Utf16BE:
		return "", errors.New("utf16BE support not yet implemented")

	case Utf8:
		return string(buf), nil

	default:
		return "", fmt.Errorf("text encoding (%v) not supported", enc)
	}
}

func (enc TextEncoding) zlen() int {
	if enc == Utf8 || enc == Iso88591 {
		return 1
	}
	if enc == Utf16 || enc == Utf16BE {
		return 2
	}
	panic(fmt.Sprintf("zlen() undefined for invalid TextEncoding (%x)", enc))
}

func toUint16s(buf []byte) ([]uint16, error) {
	if len(buf)%2 != 0 {
		return nil, fmt.Errorf("%d bytes in buffer (even number required)", len(buf))
	}
	if len(buf) == 0 {
		return []uint16{}, nil
	}

	lebom := []byte{0xff, 0xfe}
	bebom := []byte{0xfe, 0xff}
	if bytes.Equal(buf[0:2], bebom) {
		buf = buf[2:]
		for i := 0; i < len(buf)/2; i++ {
			tb := buf[i*2]
			buf[i*2] = buf[i*2+2]
			buf[i*2+2] = tb
		}
	}

	if bytes.Equal(buf[0:2], lebom) {
		buf = buf[2:]
	}

	chars := []uint16{}
	for i := 0; i < len(buf)/2; i++ {
		chars = append(chars, binary.LittleEndian.Uint16(buf[i*2:i*2+2]))
	}
	return chars, nil
}
