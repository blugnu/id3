package id3v2

import (
	"bytes"
	"testing"

	"github.com/blugnu/tags/internal/testdata"
)

func Test_header_read(t *testing.T) {

	data, err := testdata.Asset("tagged/sample.id3v22.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	header := &header{}

	src := bytes.NewReader(data)
	err = header.read(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("identifies valid header", func(t *testing.T) {
		wanted := true
		got := header.isValidHeader()
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
}
