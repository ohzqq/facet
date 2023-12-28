package facet

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func TestRoaringTerms(t *testing.T) {
	f := idx.GetFacet("tags")
	term := f.GetItem("abo")
	r := term.Bitmap()
	if len(r.ToArray()) != 416 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 416)
	}
}

func TestItemsList(t *testing.T) {
	f := idx.GetFacet("tags")
	if f.Len() != len(f.Items) {
		t.Errorf("got %d, expected %d\n", f.Len(), len(f.Items))
	}
}

func TestFuzzyFindItem(t *testing.T) {
	f := idx.GetFacet("tags")
	m := f.FuzzyFindItem("holiday")
	if len(m) != 5 {
		t.Errorf("got %d, expected 5", len(m))
	}
	//for _, i := range m {
	//fmt.Printf("%#v\n", i.Match)
	//}
}

func TestRoaringFilter(t *testing.T) {
	abo := getRoaringAbo(t)
	dnr := getRoaringDnr(t)

	or := roaring.ParOr(4, abo, dnr)
	orC := len(or.ToArray())
	if orC != 2269 {
		t.Errorf("got %d, expected %d\n", orC, 2269)
	}

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
	q, err := ParseFilters(vals)
	if err != nil {
		t.Error(err)
	}
	testFilters(q)
}

func testFilters(q url.Values) {
	items := idx.Filter(q)
	fmt.Printf("%+v\n", len(items.Data))

	//for _, item := range items.Data {
	//  fmt.Printf("%+v\n", item)
	//}
}

func getRoaringAbo(t *testing.T) *roaring.Bitmap {
	f := idx.GetFacet("tags")
	term := f.GetItem("abo")
	r := term.Bitmap()
	if len(r.ToArray()) != 416 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 416)
	}
	return r
}

func getRoaringDnr(t *testing.T) *roaring.Bitmap {
	f := idx.GetFacet("tags")
	term := f.GetItem("dnr")
	r := term.Bitmap()
	if len(r.ToArray()) != 2237 {
		t.Errorf("got %d, expected %d\n", len(r.ToArray()), 2237)
	}
	return r
}
