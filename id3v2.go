package tags

import (
	"io"
)

type Id3v2Tag struct {
	Version           Id3TagVersion
	Unsynchronisation bool
	ExtendedHeader    bool
	Experimental      bool
	Size              uint
}

func readId3v2(src io.ReadSeeker) (*Id3v2Tag, error) {
	return nil, nil
}
