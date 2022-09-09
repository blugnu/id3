package id3v2

import (
	"bytes"
	"io"

	id3reader "github.com/blugnu/tags/internal/id3/reader"
)

type Footer struct {
}

func ReadFooter(src io.ReadSeeker) (*Footer, error) {
	reader := id3reader.NewReader(src)
	sig, err := reader.ReadBytes(3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(sig, []byte(id3v2FooterSIG)) {
		return nil, nil
	}

	return &Footer{}, nil
}
