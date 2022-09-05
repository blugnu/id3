package id3v2

import (
	"errors"
	"io"

	"github.com/blugnu/tags/id3"
)

type Tag struct {
	Version           id3.TagVersion
	Unsynchronisation bool
	ExtendedHeader    bool
	Experimental      bool
	FooterPresent     bool
	Size              uint32
}

func ReadTag(src io.ReadSeeker) (*Tag, error) {
	tag := &Tag{}
	err := tag.read(src)
	if err != nil {
		if errors.Is(err, id3.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}

	return tag, nil
}

func (tag *Tag) read(src io.ReadSeeker) error {

	// Get current seek position
	opos, _ := src.Seek(0, io.SeekCurrent)

	h := header{}
	err := h.read(src)

	// If we didn't get a read error but failed to find a valid
	// header, then this is a NoTag error
	if err == nil && !h.isValidHeader() {
		err = id3.NoTag{AtPos: opos}
	}

	// Any error at this point means we either failed to read the
	// tag or the tag itself was not valid, so we restore the seek
	// position before returning any error.
	if err != nil {
		src.Seek(opos, io.SeekCurrent)
		return err
	}

	tag.Version = h.getVersion()
	tag.Size = h.tagSize

	// We don't need to reset the seek position for an unsupported version
	// of an otherwise valid v2 tag, we just don't know how to interpret the
	// tag itself.
	//
	// We will ignore it (skip over it) but there may be a further tag to
	// come which is supported.
	if tag.Version == id3.Id3vUnknown {
		return id3.UnsupportedVersionError{TagVersion: 2, Major: h.version.major, Revision: h.version.revision}
	}

	h.getFlags(
		&tag.Unsynchronisation,
		&tag.ExtendedHeader,
		&tag.Experimental,
		&tag.FooterPresent,
	)

	return nil
}
