package frame

import "bytes"

type Frame struct {
	ID                      string
	Size                    int
	PreserveWhenTagAltered  bool // 2.3.0 + 2.4.0
	PreserveWhenFileAltered bool // 2.3.0 + 2.4.0
	IsReadOnly              bool // 2.3.0 + 2.4.0
	IsCompressed            bool // 2.3.0 + 2.4.0
	IsEncrypted             bool // 2.3.0 + 2.4.0
	IsGrouped               bool // 2.3.0 + 2.4.0
	IsUnsynchronised        bool // 2.4.0
	HasDataLength           bool // 2.4.0

	TextEncoding *TextEncoding // used by text frames, otherwise nil
	LanguageCode *string       // used by comment frames
	Text         *string       // used by text frames, otherwise nil
	Description  *string       // used by user-defined text frames, otherwise nil

	UnknownData []byte // used to preserve data for otherwise unknown frame-types, otherwise nil
}

func IsValidId(id []byte) bool {
	zeroes := make([]byte, len(id))
	if bytes.Equal(id, zeroes) {
		return false
	}

	for _, b := range id {
		if (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
			continue
		}
		return false
	}

	return true
}
