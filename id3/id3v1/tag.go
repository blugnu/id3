package id3v1

import (
	"fmt"

	"github.com/blugnu/tags/id3"
)

const TagSize = 128

type Tag struct {
	Version     id3.TagVersion
	Title       string
	Artist      string
	Album       string
	Year        int
	Comment     string
	TrackNumber int
	Genre       id3.Genre
}

const Album = id3.TALB
const Title = id3.TIT2
const Artist = id3.TPE1
const Genre = id3.TCON
const Track = id3.TRCK
const Year = id3.TYER
const Comment = id3.COMM

func (tag *Tag) Get(key id3.FrameKey) string {
	switch key {

	case id3.TALB:
		return tag.Album

	case id3.TIT2:
		return tag.Title

	case id3.TPE1:
		return tag.Artist

	case id3.TCON:
		return tag.Genre.String()

	case id3.TRCK:
		return fmt.Sprintf("%d", tag.TrackNumber)

	case id3.COMM:
		return tag.Comment

	case id3.TYER:
		return fmt.Sprintf("%d", tag.Year)

	default:
		return ""
	}
}
