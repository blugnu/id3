package id3v2

import "fmt"

func (frame *framereader) readCommentFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}
	lang, err := frame.readLanguageCode()
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	d, dlen, err := frame.readStringz(enc)
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}
	c, err := frame.readString(enc, frame.Frame.Size-4-dlen) // TextEncoding + Language code = 4 bytes
	if err != nil {
		return fmt.Errorf("readComment: %w", err)
	}

	frame.Data = Comment{
		LanguageCode: lang,
		Description:  d,
		Comment:      c,
	}

	return nil
}
