package id3v2

import (
	"fmt"
	"strings"

	"github.com/blugnu/tags/id3"
)

// signature type for funcs returning frame size
type ReadFrameSizeFunc func() (uint32, error)

type framereader struct {
	id3.Reader        // ref to an id3 byte-level reader
	*Tag              // ref to the tag for which the reader is reading frames
	ReadFrameSizeFunc // ref to the function used to read frame size from the frame header
	*Frame            // a ref to the single frame obtained by an individual call to reader.readFrame()
}

func (frame *framereader) readFrame() error {
	var id string
	var size int
	var flags uint16

	err := frame.readHeader(&id, &size, &flags)
	if err != nil {
		return err
	}

	frame.Frame = &Frame{
		ID:    id,
		Size:  size,
		Flags: parseFlags(frame.Tag.Version, flags),
	}

	switch frame.ID {

	case "PIC":
	case "APIC":
		err = frame.readPictureFrame()

	case "COM":
	case "COMM":
		err = frame.readCommentFrame()

	case "TXX":
	case "TXXX":
		err = frame.readUserDefinedTextFrame()

	default:
		if strings.HasPrefix(frame.ID, "T") {
			err = frame.readTextFrame()
		} else {
			err = frame.readUnknownFrame()
		}
	}
	if err != nil {
		frame.Frame = nil
		return err
	}

	return nil
}

func (reader *framereader) readLanguageCode() (string, error) {
	buf, err := reader.ReadBytes(3)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (reader *framereader) readTextEncoding() (TextEncoding, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return UnknownTextEncoding, err
	}
	enc := TextEncoding(b)
	if !enc.isValid() {
		return UnknownTextEncoding, fmt.Errorf("invalid TextEncoding (%x)", b)
	}
	return enc, nil
}

func (reader *framereader) readString(enc TextEncoding, strlen int) (string, error) {
	buf, err := reader.ReadBytes(strlen)
	if err != nil {
		return "", fmt.Errorf("ReadString: %w", err)
	}

	s, err := enc.Decode(buf)
	if err != nil {
		return "", fmt.Errorf("ReadString: %w", err)
	}

	return s, nil
}

func (reader *framereader) readStringz(enc TextEncoding) (string, int, error) {
	buf, err := reader.ReadBytez(zlen[enc])
	if err != nil {
		return "", 0, fmt.Errorf("ReadStringz: %w", err)
	}

	s, err := enc.Decode(buf)
	if err != nil {
		return "", len(buf) + zlen[enc], fmt.Errorf("ReadStringz: %w", err)
	}

	return s, len(buf) + zlen[enc], nil
}
