package tags

import (
	"fmt"
	"io"

	"github.com/blugnu/tags/id3/id3v1"
	"github.com/blugnu/tags/id3/id3v2"
)

type AudioData struct {
	DataStart int64
	DataSize  int64
}
type Metadata struct {
	Id3v1 *id3v1.Tag
	Id3v2 []*id3v2.Tag
	Audio AudioData
}

func LoadFrom(src io.ReadSeeker) (*Metadata, error) {
	var err error

	filesize, _ := src.Seek(0, io.SeekEnd)

	md := &Metadata{
		Id3v2: []*id3v2.Tag{},
		Audio: AudioData{
			DataStart: 0,
			DataSize:  filesize,
		},
	}

	// Read any ID3v1 tag (these are always located at the end of the
	// file which the id3v1 reader takes care of)
	md.Id3v1, err = id3v1.ReadTag(src)
	if err != nil {
		return nil, fmt.Errorf("id3v1: %w", err)
	}

	// We found a v1 tag, so adjust the audio data size to account for it
	md.Audio.DataSize -= id3v1.TagSize

	// Now reposition at the start of the file and read any ID3v2 tags,
	// updating the audio data start position as we go (audio data follows
	// immediately after any id3v2 tags at the start of the file)
	src.Seek(0, io.SeekStart)
	for {
		md.Audio.DataStart, _ = src.Seek(0, io.SeekCurrent)

		tag, err := id3v2.ReadTag(src)
		if err != nil {
			md.Audio.DataStart = 0
			return nil, fmt.Errorf("id3v2: %w", err)
		}
		if tag == nil {
			break
		}

		md.Id3v2 = append(md.Id3v2, tag)
	}

	// Update the audio data size to reflect any change in the determined
	// audio data start position
	md.Audio.DataSize -= md.Audio.DataStart

	// TODO: Read any id3v2 tags located at the end of the file, updating
	// the audio data size to reflect any

	return md, nil
}
