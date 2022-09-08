package id3v2

import (
	"fmt"
)

func (frame *framereader) readText() error {
	if err := frame.readTextEncoding(); err != nil {
		return err
	}

	buf, err := frame.ReadBytes(frame.Frame.Size - 1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}

	s, err := frame.DecodeString(buf)
	if err != nil {
		return fmt.Errorf("readText: %w", err)
	}

	frame.Text = &s

	return nil
}

func (frame *framereader) readUserDefinedText() error {
	if err := frame.readTextEncoding(); err != nil {
		return err
	}

	var desc string
	var value string

	err := frame.ReadSzAndString(&desc, &value, frame.Frame.Size-1)
	if err != nil {
		return fmt.Errorf("readUserDefinedText: %w", err)
	}

	frame.Description = &desc
	frame.Text = &value

	return nil
}
