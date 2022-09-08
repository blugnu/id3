package id3v2

import (
	"bytes"
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
		return err
	}

	frame.Text = &s

	return nil
}

func (frame *framereader) readUserDefinedText() error {
	if err := frame.readTextEncoding(); err != nil {
		return err
	}
	buf, err := frame.ReadBytes(frame.Frame.Size - 1) // TextEncoding = 1 byte
	if err != nil {
		return err
	}

	el := bytes.Split(buf, frame.TextEncoding.Terminator())

	desc, err := frame.DecodeString(el[0])
	if err != nil {
		return err
	}
	comment, err := frame.DecodeString(el[1])
	if err != nil {
		return err
	}

	frame.Description = &desc
	frame.Text = &comment

	return nil
}
