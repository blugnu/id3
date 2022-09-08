package id3v2

import "fmt"

func (frame *framereader) readComment() error {
	if err := frame.readTextEncoding(); err != nil {
		return fmt.Errorf("readComment: %w", err)
	}
	if err := frame.readLanguageCode(); err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	var desc string
	var comment string

	err := frame.ReadSzAndString(&desc, &comment, frame.Frame.Size-4) // TextEncoding + Language code = 4 bytes
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	frame.Description = &desc
	frame.Text = &comment

	return nil
}
