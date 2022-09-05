package id3

type TagPlacement byte
type TagVersion byte

const (
	Id3Header TagPlacement = iota
	Id3Footer
)

const (
	Id3vUnknown TagVersion = iota
	Id3v1
	Id3v11
	Id3v22
	Id3v23
	Id3v24
)
