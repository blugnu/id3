package frame

type TextEncoding byte

const (
	Iso88591 TextEncoding = iota
	Utf16
	Utf16BE
	Utf8
)

func TextEncodingFromByte(b byte) *TextEncoding {
	enc := TextEncoding(b)
	if enc > Utf8 {
		return nil
	}
	return &enc
}

var NullTerm = map[TextEncoding][]byte{
	Iso88591: {0x00},
	Utf8:     {0x00},
	Utf16:    {0x00, 0x00},
	Utf16BE:  {0x00, 0x00},
}
