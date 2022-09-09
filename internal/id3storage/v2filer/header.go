package v2filer

import (
	"github.com/blugnu/tags/id3"
)

// Applies to all (current) Id3v2 revisions (2.2.0 / 2.3.0 / 2.4.0)
var headerflags = struct {
	compression       byte // 2.2.0 only
	unsynchronisation byte // 2.3.0 + 2.4.0
	extendedHeader    byte // 2.3.0 + 2.4.0
	experimental      byte // 2.3.0 + 2.4.0
	footer            byte // 2.3.0 + 2.4.0
}{
	compression:       0x40, // 2.2.0 only
	unsynchronisation: 0x80, // 2.3.0 + 2.4.0
	extendedHeader:    0x40, // 2.3.0 + 2.4.0
	experimental:      0x20, // 2.3.0 + 2.4.0
	footer:            0x10, // 2.3.0 + 2.4.0
}

// Extended Header flags
// NOTE: 2.2.0 does not support any extended header
type extendedheaderflagsv230 struct {
	crc uint16
}
type extendedheaderflagsv240 struct {
	update       byte
	crc          byte
	restrictions byte
}

var extendedheader = struct {
	v230flag extendedheaderflagsv230
	v240flag extendedheaderflagsv240
}{
	v230flag: extendedheaderflagsv230{
		crc: 0x8000,
	},
	v240flag: extendedheaderflagsv240{
		update:       0x40,
		crc:          0x20,
		restrictions: 0x10,
	},
}

var v240flaginfo = map[byte]struct {
	datalen byte
	name    string
}{
	extendedheader.v240flag.update:       {datalen: 0, name: "update"},
	extendedheader.v240flag.crc:          {datalen: 5, name: "crc"},
	extendedheader.v240flag.restrictions: {datalen: 1, name: "restrictions"},
}

// Maps an id3v2 major version to the corresponding TagVersion const
var tagVersion = map[byte]id3.TagVersion{
	0x02: id3.Id3v22,
	0x03: id3.Id3v23,
	0x04: id3.Id3v24,
}

const tagHeaderSIG = "ID3"
const tagFooterSIG = "3DI"
const tagHeaderSize = 10
