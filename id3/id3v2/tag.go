package id3v2

import (
	"bytes"
	"errors"
	"io"

	"github.com/blugnu/tags/id3"
	"github.com/blugnu/tags/id3/id3v2/id3v23"
	"github.com/blugnu/tags/id3/id3v2/id3v24"
	"github.com/blugnu/tags/internal/reader"
)

type Tag struct {
	Version           id3.TagVersion
	Location          int64 // the location of the tag in the source data
	Size              int   // the total size of the tag, including any extended header, frame data and padding (but not the initial header or any footer)
	IsExperimental    bool  // indicates that the tag is experimental (not the version, the tag itself)
	IsUnsynchronised  bool  // this called "Unsynchronisation" in the docs, but indicates whether unsychronisation has been applied
	HasExtendedHeader bool  // indicates whether an extended header is present
	HasFooter         bool  // indicates the presence of a footer following the frame data and any padding
	ExtendedHeader    interface{}
	Frames            []interface{}
}

func ReadTag(src io.ReadSeeker) (*Tag, error) {
	tag := &Tag{}
	if err := tag.readHeader(src); err != nil {
		if errors.Is(err, id3.NoTag{}) {
			return nil, nil
		}
		return nil, err
	}

	tagdata := make([]byte, tag.Size)
	if n, err := src.Read(tagdata); err != nil || n < tag.Size {
		return nil, err
	}

	// if tag.IsUnsynchronised {
	// TODO: apply de-unsynchronisation to tagdata before reading any frames
	// }

	reader := &framereader{reader.New(bytes.NewReader(tagdata)), tag, nil}
	for {
		err := reader.readFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tag.Frames = append(tag.Frames, reader.Frame)
	}

	// TODO: account/adjust for padding, size, current seek pos etc ?

	return tag, nil
}

func (tag *Tag) readHeader(src io.ReadSeeker) error {

	// Get current seek position
	tag.Location, _ = src.Seek(0, io.SeekCurrent)

	h := header{}
	err := h.read(src)

	// These may be "unknown" and zero for invalid or unsupported tags
	// but that's a valid outcome in that case
	tag.Version = h.getVersion()
	tag.Size = int(h.tagSize)

	// If we didn't get a read error but failed to find a valid
	// header, then there is no tag here
	if err == nil && !h.isValidHeader() {
		err = id3.NoTag{AtPos: tag.Location}
	}

	// Any error at this point means we either failed to read the
	// a valid tag or the tag has features which are not currently
	// supported.
	if err != nil {
		// If not supported, return the error and let the caller
		// decide what to do with the unsupported tag.  If the tag has
		// non-zero size, position the seek point at the end of the tag
		if errors.Is(err, id3.UnsupportedTag{}) {
			src.Seek(tag.Location+int64(tag.Size), io.SeekCurrent)
			return err
		}
		// Otherwise, position the seek point back at the beginning
		// of what we thought was going to be a tag but wasn't
		src.Seek(tag.Location, io.SeekCurrent)
		return err
	}

	h.getFlags(
		&tag.IsUnsynchronised,
		&tag.HasExtendedHeader,
		&tag.IsExperimental,
		&tag.HasFooter,
	)

	if tag.HasExtendedHeader {
		switch tag.Version {
		case id3.Id3v23:
			tag.ExtendedHeader, err = id3v23.ReadExtendedHeader(src)
		case id3.Id3v24:
			tag.ExtendedHeader, err = id3v24.ReadExtendedHeader(src)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
