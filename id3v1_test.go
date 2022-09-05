package tags

import (
	"bytes"
	"testing"

	"github.com/blugnu/tags/internal/testdata"
)

func Test_Id3v1(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v11.mp3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	seeker := bytes.NewReader(data)
	result, err := LoadFrom(seeker)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("has expected artist", func(t *testing.T) {
		wanted := "Test Artist"
		got := result.Id3v1.Artist
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected album", func(t *testing.T) {
		wanted := "Test Album"
		got := result.Id3v1.Album
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected track title", func(t *testing.T) {
		wanted := "Test Title"
		got := result.Id3v1.Title
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected comment", func(t *testing.T) {
		wanted := "Test Comment"
		got := result.Id3v1.Comment
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected track number", func(t *testing.T) {
		wanted := 3
		got := result.Id3v1.TrackNumber
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected year", func(t *testing.T) {
		wanted := 2000
		got := result.Id3v1.Year
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected genre id", func(t *testing.T) {
		wanted := Jazz
		got := result.Id3v1.Genre
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("has expected genre name", func(t *testing.T) {
		wanted := Id3GenreName[Jazz]
		got := Id3GenreName[result.Id3v1.Genre]
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}
