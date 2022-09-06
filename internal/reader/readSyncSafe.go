package reader

import (
	"encoding/binary"
)

// Reads and Decodes a 4-byte sync-safe Uint32 value
func (r *reader) ReadSyncSafeUint32() (uint32, error) {
	buf, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	// Sync-safe ints are stored in MSB order, so on a Little Endian
	// platform we need to reverse the bytes
	if byteOrder == binary.LittleEndian {
		buf = reverse(buf)
	}

	// Now can unsync the value
	return r.UnsyncUint32(buf)
}
