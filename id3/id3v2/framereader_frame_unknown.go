package id3v2

func (frame *framereader) readUnknownFrame() error {
	data, err := frame.ReadBytes(frame.Frame.Size)
	if err != nil {
		return err
	}

	frame.Data = data

	return nil

}
