package id3v2

import (
	"io"

	"github.com/blugnu/tags/id3"
	id3reader "github.com/blugnu/tags/internal/id3/reader"
)

func (tag *tagreader) readFrames(framedata []byte) error {

	// careful.. this is cute...
	//
	// when hooking up the func in the framereader to be used for ReadFrameSize,
	// we need to make sure we hookup the func provided by the framereader itself.
	//
	// so we initially get a framereader with NO ReadFrameSize func set ...
	reader := &framereader{
		id3reader.NewBytesReader(framedata),
		tag.Tag,
		nil,
		nil,
	}
	// .. and now we can give the framereader a ReadFrameSize func ref to one
	// of the funcs it provides itself!
	reader.ReadFrameSizeFunc = map[id3.TagVersion]ReadFrameSizeFunc{
		id3.Id3v22: reader.ReadUint24,
		id3.Id3v23: reader.ReadUint32,
		id3.Id3v24: reader.ReadSyncSafeUint32,
	}[tag.Version]

	for {
		err := reader.readFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		tag.Frames = append(tag.Frames, reader.Frame)
	}
	return nil
}
