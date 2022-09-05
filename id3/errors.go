package id3

import "fmt"

type NoTag struct {
	AtPos int64
}

func (e NoTag) Error() string {
	return fmt.Sprintf("no tag identified at position %d", e.AtPos)
}

// func (e NoTag) Is(target error) bool {
// 	_, ok := target.(NoTag)
// 	return ok
// }

type UnsupportedVersionError struct {
	TagVersion byte
	Major      byte
	Revision   byte
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("unsupported version (%d.%d) of a v%d tag", e.Major, e.Revision, e.TagVersion)
}
