package id3v2

import (
	"strings"

	"github.com/blugnu/tags/id3"
)

type Tag struct {
	Version           id3.TagVersion
	Location          int64  // the location of the tag in the source data
	Size              uint32 // the total size of the tag, including any extended header, frame data and padding (but not the initial header or any footer)
	IsExperimental    bool   // indicates that the tag is experimental (not the version, the tag itself)
	IsUnsynchronised  bool   // this called "Unsynchronisation" in the docs, but indicates whether unsychronisation has been applied
	HasExtendedHeader bool   // indicates whether an extended header is present
	HasFooter         bool   // indicates the presence of a footer following the frame data and any padding
	// extended header information
	IsCompressed bool
	IsUpdate     bool
	Padding      uint32
	CRC          uint32
	Restrictions byte
	// frames
	Frames []*Frame
	// raw tag data, header, including unparsed frame data, padding etc
	// (used to preserve unsupported tags, otherwise nil)
	RawData []byte
}

func (tag *Tag) Find(id string) *Frame {
	id = strings.ToUpper(id)

	for _, frame := range tag.Frames {
		if frame.ID == id {
			return frame
		}
	}
	return nil
}
