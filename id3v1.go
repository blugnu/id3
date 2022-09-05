package tags

import (
	"fmt"
	"io"
	"strings"
)

type Id3TagVersion byte

const (
	Id3v1 Id3TagVersion = iota
	Id3v11
	Id3v22
	Id3v23
	Id3v24
)

type Id3v1Tag struct {
	Version     Id3TagVersion
	Title       string
	Artist      string
	Album       string
	Year        int
	Comment     string
	TrackNumber int
	Genre       Id3Genre
}

func readId3v1(src io.ReadSeeker) (*Id3v1Tag, error) {
	_, err := src.Seek(-128, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	if tag, err := readString(src, 3); err != nil {
		return nil, err
	} else if tag != "TAG" {
		return nil, nil
	}

	tag := &Id3v1Tag{}
	tag.Title, err = readString(src, 30)
	if err != nil {
		return nil, fmt.Errorf("tag error (title): %w", err)
	}
	tag.Artist, err = readString(src, 30)
	if err != nil {
		return nil, fmt.Errorf("tag error (artist): %w", err)
	}
	tag.Album, err = readString(src, 30)
	if err != nil {
		return nil, fmt.Errorf("tag error (album): %w", err)
	}
	tag.Year, err = readStringAsInt(src, 4)
	if err != nil {
		return nil, fmt.Errorf("tag error (year): %w", err)
	}

	// Now read the comment field.  If the final byte is non-zero and the
	// preceding byte is zero, then the tag is an Id3v11 format, with the
	// track number in that final comment byte, otherwise it's an ID3 v1
	// tag with nothing but comment

	comment, err := readBytes(src, 30)
	if err != nil {
		return nil, fmt.Errorf("tag error (comment): %w", err)
	}
	if comment[28] == 0 && comment[29] != 0 {
		tag.Version = Id3v11
		tag.TrackNumber = int(comment[29])
		comment = comment[:len(comment)-2]
	}
	tag.Comment = strings.TrimSpace(string(comment))

	genre, err := readByte(src)
	if err != nil {
		return nil, fmt.Errorf("tag error (genre): %w", err)
	}
	tag.Genre = Id3Genre(genre)

	return tag, nil
}
