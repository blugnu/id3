package id3storage

import "fmt"

type NoTag struct {
	AtPos int64
}

func (e NoTag) Error() string {
	return fmt.Sprintf("no tag identified at position %d", e.AtPos)
}

func (e NoTag) Is(target error) bool {
	// A NoTag error is a NoTag error, regardless of where in the file
	// we determined there was NoTag (i.e. the error just has to be of the
	// expected type)
	_, ok := target.(NoTag)
	return ok
}

type UnsupportedTag struct {
	Reason string
}

func (e UnsupportedTag) Error() string {
	return fmt.Sprintf("unsupported tag: %s", e.Reason)
}

func (e UnsupportedTag) Is(target error) bool {
	_, ok := target.(UnsupportedTag)
	return ok
}

type UnsupportedVersionError struct {
	TagVersion byte
	Major      byte
	Revision   byte
}

func (e UnsupportedVersionError) Error() string {
	if e.Revision > 0 {
		return fmt.Sprintf("unsupported version (%d.%d) of a v%d tag", e.Major, e.Revision, e.TagVersion)
	}
	return "unsupported version"
}

type InsufficientData struct {
	Needed int
	Have   int
	Text   string
}

func (e InsufficientData) Error() string {
	if e.Text != "" {
		return fmt.Sprintf("insufficient data (needed %d bytes, have only %d): %s", e.Needed, e.Have, e.Text)
	}
	return fmt.Sprintf("insufficient data (needed %d bytes, have only %d)", e.Needed, e.Have)
}
