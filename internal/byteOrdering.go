package id3reader

import (
	"unsafe"
)

// NOTE
//
// the code in this file is currently un-used but remains present ready for the
// day that someone is able to test this code on a big-endian architecture
// and confirm for certain whether it is needed or not

type byteOrder = byte

const (
	unknown byteOrder = iota
	littleEndian
	bigEndian
)

var endianness = func() byteOrder {
	var i int32 = 1
	ptr := unsafe.Pointer(&i)
	if *(*byte)(ptr) == 0x01 {
		return littleEndian
	}
	return bigEndian
}()

// Returns a copy of the specified byte slice with the order of the bytes reversed
func reverse(bytes []byte) []byte {
	result := make([]byte, len(bytes))
	for ix, b := range bytes {
		result[len(bytes)-ix-1] = b
	}

	return result
}
