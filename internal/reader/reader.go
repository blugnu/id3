package reader

import (
	"io"
)

type reader struct {
	io.Reader
}

type Reader interface {
	ReadByte() (byte, error)
	ReadBytes(uint) ([]byte, error)
	ReadString(uint) (string, error)
	ReadStringAsInt(n uint) (int, error)
	ReadSyncSafeUint32() (uint32, error)
	UnsyncUint32([]byte) (uint32, error)
}

func New(src io.Reader) Reader {
	return &reader{src}
}
