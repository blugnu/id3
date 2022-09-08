package id3

// Reads a uint16 from a 2-byte field
func (r *reader) ReadUint16() (uint16, error) {
	ui, err := r.readUint(2)
	if err != nil {
		return 0, err
	}
	return uint16(ui), nil
}

// Reads a uint32 from a 3-byte field
func (r *reader) ReadUint24() (uint32, error) {
	ui, err := r.readUint(3)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint32 from a 4-byte field
func (r *reader) ReadUint32() (uint32, error) {
	ui, err := r.readUint(4)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint64 from a field of n bytes (maximum 8)
func (r *reader) readUint(n int) (uint, error) {
	var result uint64

	bytes, err := r.ReadBytes(n)
	if err != nil {
		return 0, err
	}

	// id3 stipulates that ints are stored in MSB order and the bit-shifting below
	// works as intended on a little-endian architecture.  however, this is UNTESTED
	// on BIG-endian architectures, so if using this code in a big-endian world,
	// be prepared for trouble!

	for _, b := range bytes {
		result = (result << 8) | uint64(b)
	}
	return uint(result), nil
}
