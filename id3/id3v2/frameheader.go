package id3v2

import "github.com/blugnu/tags/id3"

// Frame Header flags
// NOTE: 2.2.0 does not support flags in any frame header
type _v230frmhdrflags struct {
	tagAlterPreservation  uint16
	fileAlterPreservation uint16
	readonly              uint16
	compression           uint16
	encryption            uint16
	grouping              uint16
}
type _v240frmhdrflags struct {
	tagAlterPreservation  uint16
	fileAlterPreservation uint16
	readonly              uint16
	grouping              uint16
	compression           uint16
	encryption            uint16
	unsynchronisation     uint16
	datalength            uint16
}

var frameheader = struct {
	v230flag _v230frmhdrflags
	v240flag _v240frmhdrflags
}{
	v230flag: _v230frmhdrflags{
		// https://id3.org/id3v2.3.0#Frame_header_flags
		// %abc00000 ijk00000
		tagAlterPreservation:  0x8000, // a
		fileAlterPreservation: 0x4000, // b
		readonly:              0x2000, // c
		compression:           0x0080, // i
		encryption:            0x0040, // j
		grouping:              0x0020, // k
	},
	v240flag: _v240frmhdrflags{
		// https://id3.org/id3v2.4.0-structure
		// %0abc0000 0h00kmnp
		tagAlterPreservation:  0x4000, // a
		fileAlterPreservation: 0x2000, // b
		readonly:              0x1000, // c
		grouping:              0x0040, // h
		compression:           0x0008, // k
		encryption:            0x0004, // m
		unsynchronisation:     0x0002, // n
		datalength:            0x0001, // p
	},
}

// the required length of a frame ID for each id3v2 version
var idLen = map[id3.TagVersion]int{
	id3.Id3v22: 3,
	id3.Id3v23: 4,
	id3.Id3v24: 4,
}
