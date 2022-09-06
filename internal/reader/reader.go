package reader

import (
	"io"
)

type reader struct {
	curr int64
	max  int64
	io.Reader
}

type Reader interface {
	Eof() bool
	Pos() int64
	Rewind(int)
	ReadByte() (byte, error)
	ReadBytes(int) ([]byte, error)
	ReadString(int) (string, error)
	ReadStringAsInt(n int) (int, error)
	ReadSyncSafeUint32() (uint32, error)
	ReadUint16() (uint16, error)
	ReadUint24() (uint32, error)
	ReadUint32() (uint32, error)
	UnsyncUint32([]byte) (uint32, error)
}

func New(src io.Reader) Reader {
	seeker := src.(io.ReadSeeker)
	curr, _ := seeker.Seek(0, io.SeekCurrent)
	max, _ := seeker.Seek(0, io.SeekEnd)
	seeker.Seek(curr, io.SeekStart)

	return &reader{max: max, curr: curr, Reader: src}
}

func NewLimitedReader(src io.Reader, max int64) Reader {
	seeker := src.(io.ReadSeeker)
	curr, _ := seeker.Seek(0, io.SeekCurrent)
	seeker.Seek(curr, io.SeekStart)

	return &reader{max: max, curr: curr, Reader: src}
}

func (r *reader) Eof() bool {
	return r.curr == r.max
}

func (r *reader) Pos() int64 {
	return r.curr
}

func (r *reader) Rewind(n int) {
	seeker := r.Reader.(io.ReadSeeker)
	seeker.Seek(int64(-n), io.SeekCurrent)
	r.curr -= int64(n)
}
