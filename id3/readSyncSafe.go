package id3

import "fmt"

// Reads and Decodes a 4-byte sync-safe int
func (r *reader) ReadSyncSafeUint32() (uint32, error) {
	buf, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	// id3 stipulates that ints are stored in MSB order and the bit-shifting that is
	// performed when unsync-ing works as intended on a little-endian architecture.
	// this is UNTESTED on BIG-endian architectures, so if using this code in a
	// big-endian world, be prepared for trouble!

	v, err := r.UnsyncUint32(buf)
	if err != nil {
		return 0, fmt.Errorf("ReadSyncSafeUint32: %w", err)
	}

	return v, nil
}
