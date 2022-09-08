package id3

import (
	"bytes"
	"io"
)

type reader struct {
	curr int64
	max  int64
	io.Reader
}

type Reader interface {
	//	Eof() bool
	Pos() int64
	Seek(int64, int) (int64, error)
	//	Rewind(int)
	ReadByte() (byte, error)
	ReadBytes(int) ([]byte, error)
	ReadBytez(terminatorBytes []byte) ([]byte, error)
	ReadString(int) (string, error)
	ReadStringAsInt(n int) (int, error)
	ReadSyncSafeUint32() (uint32, error)
	ReadUint16() (uint16, error)
	ReadUint24() (uint32, error)
	ReadUint32() (uint32, error)
	UnsyncUint32([]byte) (uint32, error)
}

func NewBytesReader(buf []byte) Reader {
	ior := bytes.NewReader(buf)
	return NewReader(ior)
}

func NewReader(src io.Reader) Reader {
	seeker := src.(io.ReadSeeker)
	curr, _ := seeker.Seek(0, io.SeekCurrent)
	max, _ := seeker.Seek(0, io.SeekEnd)
	seeker.Seek(curr, io.SeekStart)

	return &reader{max: max, curr: curr, Reader: src}
}

func (r *reader) Eof() bool {
	return r.curr == r.max
}

func (r *reader) Pos() int64 {
	return r.curr
}

func (r *reader) Seek(pos int64, from int) (int64, error) {
	var err error
	r.curr, err = r.Reader.(io.ReadSeeker).Seek(pos, from)
	return r.curr, err
}
