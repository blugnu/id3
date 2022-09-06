package tags

import (
	"bytes"
	"testing"

	"github.com/blugnu/tags/internal/testdata"
)

func Test_LoadFrom_ID3v1Sample(t *testing.T) {

	t.Run("loads ID3v1 tag", func(t *testing.T) {
		data, err := testdata.Asset("tagged/sample.id3v11.mp3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		seeker := bytes.NewReader(data)
		result, err := LoadFrom(seeker)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
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

func Test_LoadFrom_ID3v22Sample(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v22.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seeker := bytes.NewReader(data)
	result, err := LoadFrom(seeker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("LoadFrom() returned nil")
	}

	if result.Id3v1 != nil {
		t.Fatal("Id3v1 tag is not nil")
	}

	t.Run("loads expected number of v2 tags", func(t *testing.T) {
		wanted := 1
		got := len(result.Id3v2)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})
}

func Test_LoadFrom_ID3v23Sample(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v23.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seeker := bytes.NewReader(data)
	result, err := LoadFrom(seeker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("LoadFrom() returned nil")
	}

	if result.Id3v1 != nil {
		t.Fatal("Id3v1 tag is not nil")
	}

	t.Run("loads expected number of v2 tags", func(t *testing.T) {
		wanted := 1
		got := len(result.Id3v2)
		if wanted != got {
			t.Fatalf("wanted %d, got %d", wanted, got)
		}
	})
}

func Test_LoadFrom_ID3v24Sample(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v24.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seeker := bytes.NewReader(data)
	result, err := LoadFrom(seeker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("LoadFrom() returned nil")
	}

	if result.Id3v1 != nil {
		t.Fatal("Id3v1 tag is not nil")
	}

	t.Run("loads expected number of v2 tags", func(t *testing.T) {
		wanted := 1
		got := len(result.Id3v2)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})
}
