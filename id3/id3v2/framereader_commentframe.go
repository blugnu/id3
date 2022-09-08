package id3v2

import (
	"bytes"
)

func (frame *framereader) readComment() error {
	if err := frame.readTextEncoding(); err != nil {
		return err
	}
	if err := frame.readLanguageCode(); err != nil {
		return err
	}
	buf, err := frame.ReadBytes(frame.Frame.Size - 4) // TextEncoding + Language code = 4 bytes
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
