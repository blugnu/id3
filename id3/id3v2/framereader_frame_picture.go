package id3v2

import (
	"fmt"
	"strings"
)

func (frame *framereader) readPictureFrame() error {
	enc, err := frame.readTextEncoding()
	if err != nil {
		return err
	}

	mime, nmimebytes, err := frame.readStringz(Iso88591) //
	if err != nil {
		return err
	}

	b, err := frame.ReadByte()
	if err != nil {
		return err
	}
	if b > byte(maxPictureType) {
		return fmt.Errorf("unsupported picture type (%v)", b)
	}
	pictureType := PictureType(b)

	description, ndescbytes, err := frame.readStringz(enc)
	if err != nil {
		return err
	}

	data, err := frame.ReadBytes(frame.Frame.Size - (nmimebytes + ndescbytes + 2)) // +2 = 1 each for text encoding and picture type
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

	frame.Data = Picture{
		MimeType:    mime,
		PictureType: pictureType,
		Description: description,
		Data:        data,
	}

	return nil
}
