package mp3

import "github.com/blugnu/tags/id3/id3v2"

type UnsupportedTag struct {
	*id3v2.Tag
	error
}

type audiodata struct {
	location int64
	size     int64
	bytes    []byte
}

var noaudio = &audiodata{
	location: -1,
	size:     -1,
	bytes:    nil,
}
