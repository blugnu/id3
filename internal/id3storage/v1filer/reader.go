package v1filer

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v1"
	"github.com/blugnu/tags/internal/id3storage"
)

const TagSize = 128
const SIG = "TAG"

type reader struct {
	io.ReadSeeker
	*id3v1.Tag
}

func ReadTag(src io.ReadSeeker) (*id3v1.Tag, error) {
	reader := &reader{src, nil}
	err := reader.read()
	if err != nil {
		if errors.Is(err, id3storage.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}
	return reader.Tag, nil
}

func (tag *reader) read() error {
	_, err := tag.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return err
	}

	if sig, err := tag.readV1String(3); err != nil {
		return err
	} else if sig != SIG {
		return id3storage.NoTag{AtPos: 0}
	}

	tag.Tag = &id3v1.Tag{}
	tag.Version = id3.Id3v1

	tag.Title, err = tag.readV1String(30)
	if err != nil {
		return fmt.Errorf("title: %w", err)
	}
	tag.Artist, err = tag.readV1String(30)
	if err != nil {
		return fmt.Errorf("artist: %w", err)
	}
	tag.Album, err = tag.readV1String(30)
	if err != nil {
		return fmt.Errorf("album: %w", err)
	}

	syear, err := tag.readV1String(4)
	if err != nil {
		return fmt.Errorf("year: %w", err)
	}
	year, err := strconv.ParseUint(syear, 10, 32)
	if err != nil {
		return fmt.Errorf("year (%s): %w", syear, err)
	}
	tag.Year = int(year)

	// Now read the comment field.  If the final byte is non-zero and the
	// preceding byte is zero, then the tag is an Id3v11 format, with the
	// track number in that final comment byte, otherwise it's an ID3 v1
	// tag with nothing but comment

	comment := make([]byte, 30)
	n, err := tag.Read(comment)
	if err != nil {
		return fmt.Errorf("comment: %w", err)
	}
	if n < 30 {
		return id3storage.InsufficientData{Needed: 30, Have: n, Text: "comment"}
	}
	if comment[28] == 0 && comment[29] != 0 {
		tag.Version = id3.Id3v11
		tag.TrackNumber = int(comment[29])
		comment = comment[:len(comment)-2]
	}
	tag.Comment = strings.TrimSpace(string(comment))

	genre := make([]byte, 1)
	n, err = tag.Read(genre)
	if err != nil {
		return fmt.Errorf("genre: %w", err)
	}
	if n == 0 {
		return id3storage.InsufficientData{Needed: 1, Have: n, Text: "genre"}
	}
	tag.Genre = id3.Genre(genre[0])

	// TODO: check for the presence of an extended tag (TAG+ @ -227)

	return nil
}
