package id3v2

import (
	"errors"
)

type Picture struct {
	MimeType    string
	PictureType PictureType
	Description string
	Data        []byte
}

type Frame struct {
	ID                      string
	Size                    int   // size of frame, excluding header (always 6 bytes (2.2.0) or 10 bytes (2.3.0 and 2.4.0))
	Location                int64 // location (in the file) of the frame
	PreserveWhenTagAltered  bool  // 2.3.0 + 2.4.0
	PreserveWhenFileAltered bool  // 2.3.0 + 2.4.0
	IsReadOnly              bool  // 2.3.0 + 2.4.0
	IsCompressed            bool  // 2.3.0 + 2.4.0
	IsEncrypted             bool  // 2.3.0 + 2.4.0
	IsGrouped               bool  // 2.3.0 + 2.4.0
	IsUnsynchronised        bool  // 2.4.0
	HasDataLength           bool  // 2.4.0

	TextEncoding *TextEncoding // used by text frames, otherwise nil
	LanguageCode *string       // used by comment frames
	Text         *string       // used by text frames, otherwise nil
	Description  *string       // used by user-defined text frames, otherwise nil
	Picture      *Picture      // used only by PIC/APIC frames, otherwise nil

	UnknownData []byte // used to preserve data for otherwise unknown frame-types, otherwise nil (empty = unknown data of zero length)
}

func (frame *Frame) DecodeString(buf []byte) (string, error) {
	if frame.TextEncoding == nil {
		return "", errors.New("no text encoding")
	}
	return frame.TextEncoding.Decode(buf)
}
