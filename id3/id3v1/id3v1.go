package id3v1

import (
	"errors"
	"fmt"
	"io"
	"strings"

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

func ReadTag(src io.ReadSeeker) (*Tag, error) {
	tag := &Tag{}
	err := tag.read(src)
	if err != nil {
		if errors.Is(err, id3.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}
	return tag, nil
}

func (tag *Tag) read(src io.ReadSeeker) error {
	_, err := src.Seek(-128, io.SeekEnd)
	if err != nil {
		return err
	}

	reader := id3.NewReader(src)

	if tag, err := reader.ReadString(3); err != nil {
		return err
	} else if tag != "TAG" {
		return id3.NoTag{AtPos: 0}
	}

	tag.Version = id3.Id3v1

	tag.Title, err = reader.ReadString(30)
	if err != nil {
		return fmt.Errorf("tag error (title): %w", err)
	}
	tag.Artist, err = reader.ReadString(30)
	if err != nil {
		return fmt.Errorf("tag error (artist): %w", err)
	}
	tag.Album, err = reader.ReadString(30)
	if err != nil {
		return fmt.Errorf("tag error (album): %w", err)
	}
	tag.Year, err = reader.ReadStringAsInt(4)
	if err != nil {
		return fmt.Errorf("tag error (year): %w", err)
	}

	// Now read the comment field.  If the final byte is non-zero and the
	// preceding byte is zero, then the tag is an Id3v11 format, with the
	// track number in that final comment byte, otherwise it's an ID3 v1
	// tag with nothing but comment

	comment, err := reader.ReadBytes(30)
	if err != nil {
		return fmt.Errorf("tag error (comment): %w", err)
	}
	if comment[28] == 0 && comment[29] != 0 {
		tag.Version = id3.Id3v11
		tag.TrackNumber = int(comment[29])
		comment = comment[:len(comment)-2]
	}
	tag.Comment = strings.TrimSpace(string(comment))

	genre, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("tag error (genre): %w", err)
	}
	tag.Genre = id3.Genre(genre)

	return nil
}
