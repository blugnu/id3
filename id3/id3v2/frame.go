package id3v2

import (
	"fmt"

	"github.com/blugnu/tags/id3"
)

type Frame struct {
	Key   id3.FrameKey
	ID    string
	Size  int         // size of frame, excluding header (header size is always 6 bytes (2.2.0) or 10 bytes (2.3.0 and 2.4.0))
	Flags *FrameFlags // flags from the frame header
	Data  interface{} // the frame data (either: string, []string, Comment, Picture, PositionInSet or UserDefinedText)
}

type FrameFlags struct {
	PreserveWhenTagAltered  bool // 2.3.0 + 2.4.0
	PreserveWhenFileAltered bool // 2.3.0 + 2.4.0
	IsReadOnly              bool // 2.3.0 + 2.4.0
	IsCompressed            bool // 2.3.0 + 2.4.0
	IsEncrypted             bool // 2.3.0 + 2.4.0
	IsGrouped               bool // 2.3.0 + 2.4.0
	IsUnsynchronised        bool // 2.4.0
	HasDataLength           bool // 2.4.0
}

type Comment struct {
	LanguageCode string
	Description  string
	Comment      string
}

type Picture struct {
	MimeType    string
	PictureType PictureType
	Description string
	Data        []byte
}

type PartOfSet struct {
	ItemNo    int
	ItemCount int
}

type UserDefinedText struct {
	Description string
	Text        string
}

func (com *Comment) String() string {
	return fmt.Sprintf("%s : %s", com.Description, com.Comment)
}

func (pic *Picture) String() string {
	return fmt.Sprintf("%s (%s, %d bytes)", pic.MimeType, pic.Description, len(pic.Data))
}

func (pos *PartOfSet) String() string {
	if pos.ItemNo != -1 && pos.ItemCount != -1 {
		return fmt.Sprintf("%d of %d", pos.ItemNo, pos.ItemCount)
	}
	if pos.ItemNo != -1 {
		return fmt.Sprintf("%d", pos.ItemNo)
	}
	if pos.ItemCount != -1 {
		return fmt.Sprintf("? of %d", pos.ItemCount)
	}
	return "?"
}
