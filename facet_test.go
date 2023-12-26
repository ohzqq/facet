package facet

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func TestRoaringTerms(t *testing.T) {
	f, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	term := f.GetTerm("abo")
	r := term.Bitmap()
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
	vals := make(url.Values)
	vals.Add("tags", "abo")
	vals.Add("tags", "dnr")
	vals.Add("authors", "Alice Winters")
	vals.Add("authors", "Amy Lane")
	q, err := parseFilters(vals)
	if err != nil {
		t.Error(err)
	}
	testFilters(q)
}

func testFilters(q url.Values) {
	println(q.Encode())
	items := idx.Filter(q)
	fmt.Printf("%+v\n", len(items))
	//for _, item := range items {
	//  fmt.Printf("%+v\n", item)
	//}
}

func getRoaringAbo(t *testing.T) *roaring.Bitmap {
	f, err := idx.GetFacet("tags")
	if err != nil {
		t.Error(err)
	}
	term := f.GetTerm("abo")
	r := term.Bitmap()
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
	r := term.Bitmap()
	if len(r.ToArray()) != 2237 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 2237)
	}
	return r
}
