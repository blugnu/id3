package id3v2

import "testing"

func TestIsValid(t *testing.T) {
	for _, testcase := range []struct {
		name   string
		input  byte
		result bool
	}{
		{"iso8859-1", 0, true},
		{"utf16", 1, true},
		{"utf16be", 2, true},
		{"utf8", 3, true},
		{"0x04", 0x04, false},
		{"0xff", 0xff, false},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			wanted := testcase.result
			got := TextEncoding(testcase.input).isValid()
			if wanted != got {
				t.Errorf("wanted %v, got %v", wanted, got)
			}
		})
	}
}
