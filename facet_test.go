package facet

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func TestRoaringBitmap(t *testing.T) {
	r := idx.Roar()
	if len(r.ToArray()) != 7174 {
		t.Errorf("got %d, expected %d\n", r.ToArray(), 7174)
	}
}

func TestRoaringTerms(t *testing.T) {
	f, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	term := f.GetTerm("abo")
	r := term.Roar()
	if len(r.ToArray()) != 416 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 416)
	}
}

func TestRoaringFilter(t *testing.T) {
	//books := idx.Roar()
	abo := getRoaringAbo(t)
	dnr := getRoaringDnr(t)

	//or1 := roaring.And(books, abo)
	//or2 := roaring.And(books, dnr)
	//or := roaring.Or(or1, or2)
	or := roaring.ParOr(4, abo, dnr)
	orC := len(or.ToArray())
	if orC != 2269 {
		t.Errorf("got %d, expected %d\n", orC, 2269)
	}

	//and1 := roaring.And(books, abo)
	and := roaring.ParAnd(4, abo, dnr)
	andC := len(and.ToArray())
	if andC != 384 {
		t.Errorf("got %d, expected %d\n", andC, 384)
	}
}

func TestRoaringFilters(t *testing.T) {
	q := make(url.Values)
	q.Add("tags", "abo")
	q.Add("tags", "dnr")
	q.Add("authors", "Alice Winter")
	q.Add("authors", "Amy Lane")
	items := idx.Filter(q)
	fmt.Printf("%+v\n", len(items))

	//var bits []*roaring.Bitmap
	//for name, filters := range q {
	//if facet, ok := idx.FacetCfg[name]; ok {
	//for _, filter := range filters {
	//term := facet.GetTerm(filter)
	//bits = append(bits, term.Roar())
	//}
	//}
	//}

	//or := roaring.ParOr(4, bits...)
	//orC := len(or.ToArray())
	//if orC != 2269 {
	//t.Errorf("got %d, expected %d\n", orC, 2269)
	//}

	//and := roaring.ParAnd(4, bits...)
	//andC := len(and.ToArray())
	//if andC != 384 {
	//t.Errorf("got %d, expected %d\n", andC, 384)
	//}
}

func getRoaringAbo(t *testing.T) *roaring.Bitmap {
	f, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	term := f.GetTerm("abo")
	r := term.Roar()
	if len(r.ToArray()) != 416 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 416)
	}
	return r
}

func getRoaringDnr(t *testing.T) *roaring.Bitmap {
	f, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	term := f.GetTerm("dnr")
	r := term.Roar()
	if len(r.ToArray()) != 2237 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 2237)
	}
	return r
}
