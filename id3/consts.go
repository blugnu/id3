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
	Id3v11ext
	Id3v22
	Id3v23
	Id3v24
)

var tvs = map[TagVersion]string{
	Id3v1:     "v1.0",
	Id3v11:    "v1.1",
	Id3v11ext: "v1.1 (Extended)",
	Id3v22:    "v2.2.0",
	Id3v23:    "v2.3.0",
	Id3v24:    "v2.4.0",
}

func (v TagVersion) String() string {
	if s, ok := tvs[v]; !ok {
		return "<unknown"
	} else {
		return s
	}
}

