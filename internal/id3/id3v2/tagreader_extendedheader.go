package id3v2

import (
	"fmt"
	"io"

	"github.com/blugnu/tags/id3"
)

func (tag *tagreader) readExtendedHeader() error {
	if !tag.HasExtendedHeader {
		return nil
	}

	switch tag.Version {
	case id3.Id3v23:
		return tag.readExtendedHeader230()
	case id3.Id3v24:
		return tag.readExtendedHeader240()
	}

	return nil
}

func (tag *tagreader) readExtendedHeader230() error {
	if _, err := tag.ReadSyncSafeUint32(); err != nil {
		return fmt.Errorf("readExtendedHeader230: %w", err)
	}

	flags, err := tag.ReadUint16()
	if err != nil {
		return err
	}

	tag.Padding, err = tag.ReadSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readExtendedHeader230: %w", err)
	}

	// if the CRC bit is set, read the CRC data.  In a 2.3.0 extended
	// header, this is a regular uint32, not sync-safe!
	if flags&extendedheader.v230flag.crc > 0 {
		tag.CRC, err = tag.ReadUint32()
		if err != nil {
			return err
		}
	}

	return nil
}

func (tag *tagreader) readExtendedHeader240() error {
	var err error

	pos := tag.Reader.Pos()
	size, err := tag.ReadSyncSafeUint32()
	if err != nil {
		return fmt.Errorf("readExtendedHeader240: %w", err)
	}

	// Ensure we finish with the reader positioned immediately after
	// the indicated extended header data
	defer func() {
		tag.Reader.Seek(pos+int64(size), io.SeekStart)
	}()

	// Id3v2.4 future-proofs the extended header by allowing for a variable number
	// of flag bytes, but in 2.4.0 itself there is always only 1.
	//
	// We still need to read the flag length (and should check it to ensure it is the expected value)
	flagbytes, err := tag.ReadByte()
	if err != nil {
		return err
	}
	if flagbytes != 1 {
		return fmt.Errorf("found %d flag bytes (expected 1)", flagbytes)
	}

	// read the single flag byte
	flags, err := tag.ReadByte()
	if err != nil {
		return err
	}

	// interpret data for any set flags in order of the flag bits
	for _, bit := range []byte{
		extendedheader.v240flag.update,
		extendedheader.v240flag.crc,
		extendedheader.v240flag.restrictions,
	} {
		// if the required flag bit is not set then there is nothing to read
		if flags&bit == 0 {
			continue
		}
		err = tag.readFlag(bit)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reads data for the flag identified by the flag bit specified
func (tag *tagreader) readFlag(bit byte) error {
	// Every flag consists of:
	// - a byte identifying the length of any additional data for the flag
	// - any additional data bytes (or none)

	datalen, err := tag.ReadByte()
	if err != nil {
		return err
	}
	expected := v240flaginfo[bit]
	if datalen != expected.datalen {
		return fmt.Errorf(" %s flag: unexpected data length: got %d, expected %d", expected.name, datalen, expected.datalen)
	}
	data, err := tag.ReadBytes(int(datalen))
	if err != nil {
		return err
	}

	// now parse the flag data according to the flag type
	switch bit {
	// update flag - no data, it's mere presence IS the data
	case extendedheader.v240flag.update:
		tag.IsUpdate = true
		return nil

	// crc flag - has 5 bytes of data containing a 32-bit sync-safe int
	case extendedheader.v240flag.crc:
		tag.CRC, err = tag.UnsyncUint32(data)
		if err != nil {
			return err
		}
		return nil

		// restrictions flag - 1 bit of additional bits signifiying varying restrictions
	case extendedheader.v240flag.restrictions:
		// TODO: parse the bits and set restriction indicators in the tag (to be defined)
		// For now, having read the data byte is enough to "handle" this flag so we'll just
		// set the tag.Restrictions to the unparsed byte
		tag.Restrictions = data[0]
		return nil
	}

	return fmt.Errorf("flag bit not supported (%x)", bit)
}
