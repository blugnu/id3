package id3v1

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v1"
	id3reader "github.com/blugnu/tags/internal/id3/reader"
)

const SIG = "TAG"
const TagSize = 128

type reader struct {
	id3reader.Reader
	*id3v1.Tag
}

func ReadTag(src io.ReadSeeker) (*id3v1.Tag, error) {
	tag := &reader{id3reader.NewReader(src), nil}

	err := tag.read()
	if err != nil {
		if errors.Is(err, id3reader.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}
	return tag.Tag, nil
}

func (tag *reader) read() error {
	_, err := tag.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return err
	}

	if sig, err := tag.readString(3); err != nil {
		return err
	} else if sig != SIG {
		return id3reader.NoTag{AtPos: 0}
	}

	tag.Tag = &id3v1.Tag{}
	tag.Version = id3.Id3v1

	tag.Title, err = tag.readString(30)
	if err != nil {
		return fmt.Errorf("tag error (title): %w", err)
	}
	tag.Artist, err = tag.readString(30)
	if err != nil {
		return fmt.Errorf("tag error (artist): %w", err)
	}
	tag.Album, err = tag.readString(30)
	if err != nil {
		return fmt.Errorf("tag error (album): %w", err)
	}
	tag.Year, err = tag.readIntString(4)
	if err != nil {
		return fmt.Errorf("tag error (year): %w", err)
	}

	// Now read the comment field.  If the final byte is non-zero and the
	// preceding byte is zero, then the tag is an Id3v11 format, with the
	// track number in that final comment byte, otherwise it's an ID3 v1
	// tag with nothing but comment

	comment, err := tag.ReadBytes(30)
	if err != nil {
		return fmt.Errorf("tag error (comment): %w", err)
	}
	if comment[28] == 0 && comment[29] != 0 {
		tag.Version = id3.Id3v11
		tag.TrackNumber = int(comment[29])
		comment = comment[:len(comment)-2]
	}
	tag.Comment = strings.TrimSpace(string(comment))

	genre, err := tag.ReadByte()
	if err != nil {
		return fmt.Errorf("tag error (genre): %w", err)
	}
	tag.Genre = id3.Genre(genre)

	// TODO: check for the presence of an extended tag (TAG+ @ -227)

	return nil
}
