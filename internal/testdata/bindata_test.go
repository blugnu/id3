package testdata

import "testing"

func TestData(t *testing.T) {
	t.Run("provides the expected number of original samples", func(t *testing.T) {
		files, err := AssetDir("original")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		wanted := 5
		got := len(files)
		if wanted != got {
			t.Errorf("wanted %v, got %d", wanted, got)
		}
	})

	t.Run("provides the expected number of tagged samples", func(t *testing.T) {
		files, err := AssetDir("tagged")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		wanted := 10
		got := len(files)
		if wanted != got {
			t.Errorf("wanted %v, got %d", wanted, got)
		}
	})
}
