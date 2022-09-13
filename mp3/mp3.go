package mp3

import (
	"bytes"
	"fmt"
	"os"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v1"
	"github.com/blugnu/tags/id3/id3v2"
)

type mp3 struct {
	filename string
	*audiodata
	Id3v1           *id3v1.Tag
	Id3v2           []*id3v2.Tag
	UnsupportedTags []UnsupportedTag
}

func FromBytes(buf []byte) (*mp3, error) {
	seeker := bytes.NewReader(buf)
	return read(seeker)
}

func FromFile(filename string) (*mp3, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("FromFile: %w", err)
	}
	defer file.Close()

	mp3, err := read(file)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	mp3.filename = filename

	return mp3, nil
}

func (mp3 *mp3) CreateTag(ver id3.TagVersion) *id3v2.Tag {
	tag := &id3v2.Tag{
		Version: ver,
		Frames:  []*id3v2.Frame{},
	}
	mp3.Id3v2 = append(mp3.Id3v2, tag)
	return tag
}

func (mp3 *mp3) GetTag(ver id3.TagVersion) *id3v2.Tag {
	for _, tag := range mp3.Id3v2 {
		if tag.Version == ver {
			return tag
		}
	}
	return nil
}
