package id3v2

import (
	"fmt"
	"strconv"

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

func (tag *Tag) Find(key id3.FrameKey) *Frame {
	for _, frame := range tag.Frames {
		if frame.Key == key {
			return frame
		}
	}
	return nil
}

func (tag *Tag) Get(key id3.FrameKey) string {
	frame := tag.Find(key)
	if frame == nil {
		return ""
	}

	switch frame.Data.(type) {
	case *PartOfSet:
		return frame.Data.(*PartOfSet).String()

	case string:
		return frame.Data.(string)
	}

	return ""
}

func (tag *Tag) GetInt(key id3.FrameKey) int {
	var frame *Frame

	switch key {
	case id3.DiscNo, id3.NumDiscs:
		frame = tag.Find(id3.TPOS)
	case id3.TrackNo, id3.NumTracks:
		frame = tag.Find(id3.TRCK)
	default:
		frame = tag.Find(key)
	}
	if frame == nil {
		return -1
	}

	switch frame.Data.(type) {
	case *PartOfSet:
		pos := frame.Data.(*PartOfSet)
		if key == id3.NumDiscs || key == id3.NumTracks {
			return pos.ItemCount
		}
		return pos.ItemNo

	case string:
		text, ok := frame.Data.(string)
		if !ok {
			i, err := strconv.ParseInt(text, 10, 32)
			if err != nil {
				return -1
			}
			return int(i)
		}
		return -1

	default:
		return -1
	}
}

func (tag *Tag) Set(key id3.FrameKey, value string) {
	frame := tag.Find(key)
	if frame == nil {
		return
	}

	switch frame.Data.(type) {
	case string:
		frame.Data = value
		return
	}
}

func (tag *Tag) SetInt(key id3.FrameKey, value int) {
	var frame *Frame

	switch key {
	case id3.DiscNo, id3.NumDiscs:
		frame = tag.Find(id3.TPOS)
	case id3.TrackNo, id3.NumTracks:
		frame = tag.Find(id3.TRCK)
	default:
		frame = tag.Find(key)
	}
	if frame == nil {
		return
	}

	switch frame.Data.(type) {
	case *PartOfSet:
		pos := frame.Data.(*PartOfSet)
		if key == id3.NumDiscs || key == id3.NumTracks {
			pos.ItemCount = value
			return
		}
		pos.ItemNo = value
		return

	case string:
		frame.Data = fmt.Sprintf("%d", value)
		return
	}
}
