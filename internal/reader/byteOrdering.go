package reader

import (
	"unsafe"
)

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

// Returns the specified byte slice with the order of the bytes reversed
func reverse(bytes []byte) []byte {
	result := make([]byte, len(bytes))
	for ix, b := range bytes {
		result[len(bytes)-ix-1] = b
	}

	return result
}
