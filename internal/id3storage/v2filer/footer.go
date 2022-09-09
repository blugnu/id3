package v2filer

import (
	"bytes"
	"io"
)

type Footer struct {
}

func ReadFooter(src io.ReadSeeker) (*Footer, error) {
	reader := &reader{src}
	sig, err := reader.readBytes(3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(sig, []byte(tagFooterSIG)) {
		return nil, nil
	}

	return &Footer{}, nil
}
