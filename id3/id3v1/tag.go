package id3v1

import (
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
