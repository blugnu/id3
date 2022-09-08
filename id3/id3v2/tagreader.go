package id3v2

import (
	"errors"
	"io"

	"github.com/blugnu/tags/id3"
)

type tagreader struct {
	id3.Reader
	*Tag
}

func (tag *tagreader) readTag() error {

	var pos int64
	var majorver byte
	var revision byte
	var flags byte
	var size uint32

	err := tag.readHeader(&pos, &majorver, &revision, &flags, &size)
	defer func() {
		// a non-zero size means we at least THINK we read a valid header
		// and should make sure we position at the end of the tag when
		// we're done, whatever happens
		if size > 0 {
			tag.Seek(pos+tagHeaderSize+int64(size), io.SeekStart)
		}
	}()

	// if we get an UnsupportedTag error, this indicates we read a
	// header that was identified as a v2 tag ("ID3") but is of a version
	// or uses features that are not supported by this implementation
	//
	// in which case we identify it as a tag where the only known
	// information is the location, the size (excluding header) and
	// raw tag data (INCLUDING the entire header)
	if err != nil {
		if errors.Is(err, id3.UnsupportedTag{}) {
			tag.Seek(pos, io.SeekStart)
			data, err := tag.ReadBytes(int(size) + tagHeaderSize)
			if err != nil {
				return err
			}
			tag.Tag = &Tag{Location: pos, Size: size, raw: data}
		}
		return err
	}

	// otherwise we have a valid, supported tag...

	tag.Tag = &Tag{
		Version:           tagVersion[majorver],
		Size:              size,
		Location:          pos,
		IsExperimental:    flags&headerflags.unsynchronisation > 0,
		IsUnsynchronised:  flags&headerflags.extendedHeader > 0,
		HasExtendedHeader: flags&headerflags.experimental > 0,
		HasFooter:         flags&headerflags.footer > 0,
	}

	// .. which may also have an extended header

	if err := tag.readExtendedHeader(); err != nil {
		return err
	}

	// .. and will definitely have at least one frame

	framedata, err := tag.ReadBytes(int(tag.Size))
	if err != nil {
		return err
	}

	// if tag.IsUnsynchronised {
	// TODO: apply de-unsynchronisation to tagdata before reading any frames
	// }

	err = tag.readFrames(framedata)
	if err != nil {
		return err
	}

	return nil
}
