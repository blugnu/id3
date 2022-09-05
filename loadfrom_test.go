package tags

import (
	"bytes"
	"testing"

	"github.com/blugnu/tags/internal/testdata"
)

func Test_LoadFrom(t *testing.T) {

	t.Run("loads ID3v1 tag", func(t *testing.T) {
		data, err := testdata.Asset("tagged/sample.id3v11.mp3")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		seeker := bytes.NewReader(data)
		result, err := LoadFrom(seeker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("LoadFrom() returned nil")
		}

		if result.Id3v1 == nil {
			t.Fatal("Id3v1 tag is nil")
		}

		wanted := "Test Artist"
		got := result.Id3v1.Artist
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}
