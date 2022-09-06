package reader

import "errors"

// Desyncs the specified byte slice and returns the result as a uint32.
// The supplied slice may contain a maximum of 5 bytes.
func (*reader) UnsyncUint32(bytes []byte) (uint32, error) {
	var result uint32
	for _, b := range bytes {
		// If the high bit is set then this isn't a valid sync-safe integer
		if b&0x80 > 0 {
			return 0, errors.New("not a sync safe integer")
		}
		result = (result << 7) | uint32(b)
	}
	return result, nil
}
