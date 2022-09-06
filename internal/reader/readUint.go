package reader

// Reads a uint16 from a 2-byte field
func (r *reader) ReadUint16() (uint16, error) {
	ui, err := r.ReadUint(2)
	if err != nil {
		return 0, err
	}
	return uint16(ui), nil
}

// Reads a uint32 from a 3-byte field
func (r *reader) ReadUint24() (uint32, error) {
	ui, err := r.ReadUint(3)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint32 from a 4-byte field
func (r *reader) ReadUint32() (uint32, error) {
	ui, err := r.ReadUint(4)
	if err != nil {
		return 0, err
	}
	return uint32(ui), nil
}

// Reads a uint64 from a field of n bytes (maximum 8)
func (r *reader) ReadUint(n int) (uint64, error) {
	var result uint64

	bytes, err := r.ReadBytes(n)
	if err != nil {
		return 0, err
	}

	if endianness == bigEndian {
		bytes = reverse(bytes)
	}

	for _, b := range bytes {
		result = (result << 8) | uint64(b)
	}
	return result, nil
}
