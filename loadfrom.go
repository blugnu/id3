package tags

import (
	"fmt"
	"io"
)

type Metadata struct {
	Id3v1         *Id3v1Tag
	Id3v2         []*Id3v2Tag
	AudioDataPos  int64
	AudioDataSize int64
}

func LoadFrom(src io.ReadSeeker) (*Metadata, error) {
	var err error

	md := &Metadata{
		Id3v2:         []*Id3v2Tag{},
		AudioDataPos:  -1,
		AudioDataSize: -1,
	}

	// Read any ID3v1 tag (from the end of the file)
	md.Id3v1, err = readId3v1(src)
	if err != nil {
		return nil, fmt.Errorf("id3v1: %w", err)
	}

	// Read any ID3v2 tags (from the start of the file)
	src.Seek(0, io.SeekStart)
	for {
		md.AudioDataPos, _ = src.Seek(0, io.SeekCurrent)

		tag, err := readId3v2(src)
		if err != nil {
			md.AudioDataPos = -1
			return nil, fmt.Errorf("id3v2: %w", err)
		}
		if tag != nil {
			md.Id3v2 = append(md.Id3v2, tag)
			continue
		}

		break
	}

	return md, nil
}
