package facet

import (
	"testing"
)

func TestRoaringBitmap(t *testing.T) {
	r := idx.Roar()
	if len(r.ToArray()) != 7174 {
		t.Errorf("got %d, expected %d\n", r.ToArray(), 7174)
	}
}

func TestRoaringTerms(t *testing.T) {
}
