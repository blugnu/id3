package id3v2

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
	UnknownTextEncoding = 0xff
)

func (enc TextEncoding) isValid() bool {
	return enc <= Utf8
}

func (enc TextEncoding) Decode(buf []byte) (string, error) {
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

var zlen = map[TextEncoding]int{
	Iso88591: 1,
	Utf8:     1,
	Utf16:    2,
	Utf16BE:  2,
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
