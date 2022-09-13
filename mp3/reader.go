package mp3

import (
	"fmt"
	"io"

	"github.com/blugnu/tags/id3/id3v1"
	"github.com/blugnu/tags/id3/id3v2"
	"github.com/blugnu/tags/internal/id3storage/v1filer"
	"github.com/blugnu/tags/internal/id3storage/v2filer"
)

func read(src io.ReadSeeker) (*mp3, error) {
	var err error

	filesize, _ := src.Seek(0, io.SeekEnd)

	mp3 := &mp3{
		Id3v2: []*id3v2.Tag{},
		audiodata: &audiodata{
			location: 0,
			size:     filesize,
		},
	}

	// Read any ID3v1 tag (these are always located at the END of the
	// file which the id3v1 reader takes care of)
	mp3.Id3v1, err = v1filer.ReadTag(src)
	if err != nil {
		return nil, fmt.Errorf("id3v1: %w", err)
	}

	// if we found a v1 tag adjust the audio data size to account for it
	if mp3.Id3v1 != nil {
		mp3.audiodata.size -= id3v1.TagSize
	}

	// reposition at the start of the file and read any ID3v2 tags
	src.Seek(0, io.SeekStart)
	for {
		tag, err := v2filer.ReadTag(src)
		if tag != nil {
			// Update the audio data location and size to account for the tag
			mp3.audiodata.location = tag.Location + int64(tag.Size)
			mp3.audiodata.size -= int64(tag.Size)

			if err == nil {
				// We got a tag, with no errors, so add it to the
				// v2 tags and look for another one
				mp3.Id3v2 = append(mp3.Id3v2, tag)
				continue
			}
			// We got a tag, but also an error, so the tag is
			// considered unsupported
			mp3.UnsupportedTags = append(mp3.UnsupportedTags, UnsupportedTag{tag, err})
			continue
		}

		// No tag, but an error, something catastrophic has happened.  We can
		// keep any tags we may have found to this point but we cannot rely
		// on the audio positioning or size, so we clobber those before
		// returning the error.
		//
		// The information in any tags extracted so far can be used, but the
		// information cannot be written back to the source file
		if err != nil {
			mp3.audiodata = noaudio
			return mp3, fmt.Errorf("id3v2: %w", err)
		}

		// No error, but no tag... we're done looking for tags
		break
	}

	// TODO: Read any id3v2 tags located at the END of the file, updating
	// the audio data size to reflect any

	// reposition at the END of the file and check for a tag footer

	pos := int64(-10)
	for {
		src.Seek(pos, io.SeekEnd)
		footer, err := v2filer.ReadFooter(src)
		if err != nil {
			mp3.audiodata = noaudio
			return mp3, fmt.Errorf("read [footer]: %w", err)
		}
		if footer == nil {
			break
		}
		println("something's afoot!")
	}

	// load the audio data
	audioBytesRead := 0
	_, err = src.Seek(mp3.audiodata.location, io.SeekStart)
	if err == nil {
		mp3.audiodata.bytes = make([]byte, mp3.audiodata.size)
		audioBytesRead, err = src.Read(mp3.audiodata.bytes)
	}

	if err != nil || int64(audioBytesRead) < mp3.audiodata.size {
		mp3.audiodata = noaudio
	}

	return mp3, err
}
