package id3v2

import (
	"fmt"
	"strings"
)

func (frame *framereader) readPicture() error {
	if err := frame.readTextEncoding(); err != nil {
		return err
	}

	enc := Iso88591
	buf, err := frame.ReadBytez(enc.Terminator())
	if err != nil {
		return err
	}
	mime, err := enc.Decode(buf)
	if err != nil {
		return err
	}
	lenImageType := len(buf) + 1

	b, err := frame.ReadByte()
	if err != nil {
		return err
	}
	if b > byte(maxPictureType) {
		return fmt.Errorf("unsupported picture type (%v)", b)
	}
	pictureType := PictureType(b)

	enc = *frame.TextEncoding
	buf, err = frame.ReadBytez(enc.Terminator())
	if err != nil {
		return err
	}
	description, err := enc.Decode(buf)
	if err != nil {
		return err
	}
	lenDescription := len(buf) + len(enc.Terminator())

	data, err := frame.ReadBytes(frame.Frame.Size - (1 + lenImageType + 1 + lenDescription))
	if err != nil {
		return err
	}

	mime = strings.ToLower(mime)
	switch mime {
	case "gif":
		mime = "image/gif"
	case "png":
		mime = "image/png"
	case "jpg":
	case "jpeg":
		mime = "image/jpeg"
	}

	frame.Picture = &Picture{
		MimeType:    mime,
		PictureType: pictureType,
		Description: description,
		Data:        data,
	}

	return nil
}
