package reader

import (
	"encoding/binary"
	"unsafe"
)

var byteOrder = func() binary.ByteOrder {
	var i int = 0x0100
	ptr := unsafe.Pointer(&i)
	if *(*byte)(ptr) == 0x01 {
		return binary.BigEndian
	}
	return binary.LittleEndian
}()

// Returns the specified byte slice with the order of the bytes reversed
func reverse(bytes []byte) []byte {
	result := make([]byte, len(bytes))
	for ix, b := range bytes {
		result[len(bytes)-ix-1] = b
	}

	return result
}
