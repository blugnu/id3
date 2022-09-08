package id3v2

import (
	"fmt"
)

func (frame *framereader) readTextFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	buf, err := frame.ReadBytes(frame.Frame.Size - 1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}

	v, err := enc.Decode(buf)
	if err != nil {
		return fmt.Errorf("readText: %w", err)
	}

	frame.Data = v

	return nil
}

func (frame *framereader) readUserDefinedTextFrame() error {

	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	d, dlen, err := frame.readStringz(enc)
	if err != nil {
		return fmt.Errorf("readUserDefinedText: %w", err)
	}

	t, err := frame.readString(enc, frame.Frame.Size-1-dlen)
	if err != nil {
		return fmt.Errorf("readUserDefinedText: %w", err)
	}

	frame.Data = UserDefinedText{
		Description: d,
		Text:        t,
	}

	return nil
}
