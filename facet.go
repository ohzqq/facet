package facet

import (
	"net/url"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Facet struct {
	Name     string
	Terms    map[string]*Term
	Operator string
	terms    url.Values
}

type Term struct {
	Value string
	Count int
	items []uint32
}

func NewFacet(name string) *Facet {
	return &Facet{
		Name:     name,
		Operator: "or",
		Terms:    make(map[string]*Term),
		terms:    make(url.Values),
	}
}

func (f *Facet) GetTerm(term string) *Term {
	if t, ok := f.Terms[term]; ok {
		return t
	}
	return &Term{Value: term}
}

func (f *Facet) Filter(filters ...string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, filter := range filters {
		term := f.GetTerm(filter)
		bits = append(bits, term.Bitmap())
	}

	switch f.Operator {
	case "and":
		return roaring.ParAnd(4, bits...)
	default:
		return roaring.ParOr(4, bits...)
	}
}

func collectFacetValues(name string, data []map[string]any) url.Values {
	facet := make(url.Values)
	for i, item := range data {
		if terms, ok := item[name]; ok {
			for _, term := range terms.([]any) {
				facet.Add(cast.ToString(term), cast.ToString(i))
			}
		}
	}
	return facet
}

func NewTerm(name string, vals []string) *Term {
	term := &Term{
		Value: name,
		Count: len(vals),
		items: make([]uint32, len(vals)),
	}
	for i, val := range vals {
		term.items[i] = cast.ToUint32(val)
	}
	return term
}

func (t *Term) Bitmap() *roaring.Bitmap {
	return roaring.BitmapOf(t.items...)
}
