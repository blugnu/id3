package v1filer

import (
	"bytes"
	"testing"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/internal/testdata"
)

func Test_Id3v1(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v11.mp3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	src := bytes.NewReader(data)
	tag, err := ReadTag(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("has expected artist", func(t *testing.T) {
		wanted := "Test Artist"
		got := tag.Artist
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected album", func(t *testing.T) {
		wanted := "Test Album"
		got := tag.Album
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected track title", func(t *testing.T) {
		wanted := "Test Title"
		got := tag.Title
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected comment", func(t *testing.T) {
		wanted := "Test Comment"
		got := tag.Comment
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected track number", func(t *testing.T) {
		wanted := 3
		got := tag.TrackNumber
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected year", func(t *testing.T) {
		wanted := 2000
		got := tag.Year
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected genre id", func(t *testing.T) {
		wanted := id3.Jazz
		got := tag.Genre
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected genre name", func(t *testing.T) {
		wanted := id3.Jazz.String()
		got := tag.Genre.String()
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}
