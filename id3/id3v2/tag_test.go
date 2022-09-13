package id3v2

import (
	"testing"

	"github.com/blugnu/tags/id3"
)

func TestTPOS(t *testing.T) {
	tpos := &Frame{
		ID:   "TPOS",
		Key:  id3.TPOS,
		Data: &PartOfSet{ItemNo: 1, ItemCount: 2},
	}
	tag := &Tag{
		Frames: []*Frame{tpos},
	}

	t.Run("Get(TPOS))", func(t *testing.T) {
		wanted := "1 of 2"
		got := tag.Get(id3.TPOS)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("GetInt(DiscNo)", func(t *testing.T) {
		wanted := 1
		got := tag.GetInt(id3.DiscNo)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
}

func TestTRCK(t *testing.T) {
	tpos := &Frame{
		ID:   "TRCK",
		Key:  id3.TRCK,
		Data: &PartOfSet{ItemNo: 16, ItemCount: -1},
	}
	tag := &Tag{
		Frames: []*Frame{tpos},
	}

	t.Run("Get(TRCK)", func(t *testing.T) {
		wanted := "16"
		got := tag.Get(id3.TRCK)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("GetInt(TrackNo)", func(t *testing.T) {
		wanted := 16
		got := tag.GetInt(id3.TrackNo)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("GetInt(NumTracks)", func(t *testing.T) {
		wanted := -1
		got := tag.GetInt(id3.NumTracks)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("SetInt(NumTracks)", func(t *testing.T) {
		tag.SetInt(id3.NumTracks, 37)
		wanted := 37
		got := tag.GetInt(id3.NumTracks)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Get(TRACK) (with NumTracks set)", func(t *testing.T) {
		tag.SetInt(id3.NumTracks, 37)
		wanted := "16 of 37"
		got := tag.Get(id3.TRCK)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
}
