package id3v2

import (
	"fmt"
	"io"

	"github.com/blugnu/tags/id3"
)

func (header *framereader) readHeader(id *string, size *int, flags *uint16) error {
	var err error

	idbytes, err := header.ReadBytes(idLen[header.Tag.Version])
	if err != nil {
		return err
	}
	if !header.isValidId(idbytes) {
		return io.EOF
	}
	*id = string(idbytes)

	// we got a valid frame id, so now we can read the frame size...
	//ui32, err := header.ReadFrameSizeFunc(header.Reader)
	ui32, err := header.ReadFrameSizeFunc()
	if err != nil {
		return fmt.Errorf("framereader.readHeader.ReadFrameSize: %w", err)

	}
	*size = int(ui32)

	// ... for a 2.2.0 frame the header is complete...
	if header.Tag.Version == id3.Id3v22 {
		return nil
	}

	// ... but for 2.30 and 2.40 the header also contains 2 additional
	// bytes of frame flags.  We'll need to parse those later, once we
	// have a frame to parse them into; for now we just set the flag
	// bytes before returning
	*flags, err = header.ReadUint16()
	if err != nil {
		return err
	}

	return nil
}

func (reader *framereader) parseHeaderFlags(flags uint16) {
	frame := reader.Frame

	switch reader.Tag.Version {
	case id3.Id3v23:
		frame.PreserveWhenTagAltered = flags&frameheader.v230flag.tagAlterPreservation > 0
		frame.PreserveWhenFileAltered = flags&frameheader.v230flag.fileAlterPreservation > 0
		frame.IsReadOnly = flags&frameheader.v230flag.readonly > 0
		frame.IsCompressed = flags&frameheader.v230flag.compression > 0
		frame.IsEncrypted = flags&frameheader.v230flag.encryption > 0
		frame.IsGrouped = flags&frameheader.v230flag.grouping > 0

	case id3.Id3v24:
		frame.PreserveWhenTagAltered = flags&frameheader.v240flag.tagAlterPreservation > 0
		frame.PreserveWhenFileAltered = flags&frameheader.v240flag.fileAlterPreservation > 0
		frame.IsReadOnly = flags&frameheader.v240flag.readonly > 0
		frame.IsGrouped = flags&frameheader.v240flag.grouping > 0
		frame.IsCompressed = flags&frameheader.v240flag.compression > 0
		frame.IsEncrypted = flags&frameheader.v240flag.encryption > 0
		frame.IsUnsynchronised = flags&frameheader.v240flag.unsynchronisation > 0
		frame.HasDataLength = flags&frameheader.v240flag.datalength > 0
	}
}

func (tag *framereader) isValidId(id []byte) bool {
	// the id must be of the correct length
	if len(id) != idLen[tag.Version] {
		return false
	}

	// it must consist only of A..Z or 0..9 chars
	for _, b := range id {
		if (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
			continue
		}
		return false
	}

	return true
}
