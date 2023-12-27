package facet

import (
	"net/url"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Facet struct {
	Attribute string           `json:"attribute"`
	Items     map[string]*Term `json:"items,omitempty"`
	Operator  string           `json:"operator,omitempty"`
}

type Term struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
	items []uint32
}

func NewFacet(name string) *Facet {
	return &Facet{
		Attribute: name,
		Operator:  "or",
		Items:     make(map[string]*Term),
	}
}

func (f *Facet) GetTerm(term string) *Term {
	if t, ok := f.Items[term]; ok {
		return t
	}
	return NewTerm(term, []string{})
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
		Label: name,
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
