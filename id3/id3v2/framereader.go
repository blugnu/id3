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
	loc := frame.Pos()

	err := frame.readHeader(&id, &size, &flags)
	if err != nil {
		return err
	}

	frame.Frame = &Frame{
		ID:       id,
		Size:     size,
		Location: loc,
	}
	frame.parseHeaderFlags(flags)

	switch frame.ID {

	case "PIC":
	case "APIC":
		err = frame.readPicture()

	case "COM":
	case "COMM":
		err = frame.readComment()

	case "TXX":
	case "TXXX":
		err = frame.readUserDefinedText()

	default:
		if strings.HasPrefix(frame.ID, "T") {
			err = frame.readText()
		} else {
			err = frame.readUnknown()
		}
	}
	if err != nil {
		frame.Frame = nil
		return err
	}

	return nil
}

func (frame *framereader) readUnknown() error {
	data, err := frame.ReadBytes(frame.Frame.Size)
	if err != nil {
		return err
	}

	frame.UnknownData = data

	return nil

}

func (reader *framereader) readLanguageCode() error {
	buf, err := reader.ReadBytes(3)
	if err != nil {
		return err
	}
	sbuf := string(buf)
	reader.Frame.LanguageCode = &sbuf
	return nil
}

func (reader *framereader) readTextEncoding() error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	reader.Frame.TextEncoding = TextEncodingFromByte(b)
	return nil
}

func (reader *framereader) ReadSzAndString(sz *string, s *string, totalbytes int) error {

	enc := reader.Frame.TextEncoding
	terminator := enc.Terminator()

	buf, err := reader.ReadBytez(terminator)
	if err != nil {
		return fmt.Errorf("ReadSzAndString [sz]: %w", err)
	}
	*sz, err = enc.Decode(buf)
	if err != nil {
		return fmt.Errorf("ReadSzAndString [sz]: %w", err)
	}

	buf, err = reader.ReadBytes(totalbytes - (len(buf) + len(terminator)))
	if err != nil {
		return fmt.Errorf("ReadSzAndString [s]]: %w", err)
	}
	*s, err = enc.Decode(buf)
	if err != nil {
		return fmt.Errorf("ReadSzAndString [value]: %w", err)
	}

	return nil
}
