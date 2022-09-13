package id3

import "testing"

func TestFrameKeySet(t *testing.T) {
	t.Run("IsEmpty on an empty set", func(t *testing.T) {
		set := FrameKeySet{}
		wanted := true
		got := set.IsEmpty()
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
	t.Run("IsEmpty on a non-empty set", func(t *testing.T) {
		set := FrameKeySet{TALB}
		wanted := false
		got := set.IsEmpty()
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Add to an empty set", func(t *testing.T) {
		set := FrameKeySet{}
		set = set.Add(TALB)
		wanted := 1
		got := len(set)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Add duplicate to a set", func(t *testing.T) {
		set := FrameKeySet{TALB}
		set = set.Add(TALB)
		wanted := 1
		got := len(set)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Contains with an empty set", func(t *testing.T) {
		set := FrameKeySet{}
		wanted := false
		got := set.Contains(TALB)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Contains with set having key of interest", func(t *testing.T) {
		set := FrameKeySet{TALB}
		wanted := true
		got := set.Contains(TALB)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Remove from an empty set", func(t *testing.T) {
		set := FrameKeySet{}
		set = set.Remove(TALB)
		wanted := 0
		got := len(set)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("Remove key from a set", func(t *testing.T) {
		set := FrameKeySet{TALB, TCOM, APIC}
		set = set.Remove(TALB)
		wanted := 2
		got := len(set)
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

}
