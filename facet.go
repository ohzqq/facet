package facet

import (
	"net/url"

	"github.com/kelindar/bitmap"
	"github.com/samber/lo"
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
	items []string
}

func NewFacet(name string) *Facet {
	return &Facet{
		Name:     name,
		Operator: "or",
	}
}

func NewFacetVals(name string, pk string, data []map[string]any) url.Values {
	facet := make(url.Values)
	for _, item := range data {
		if terms, ok := item[name]; ok {
			for _, term := range terms.([]any) {
				facet.Add(cast.ToString(term), cast.ToString(item[pk]))
			}
		}
	}
	return facet
}

//func Intersect(vals url.Values, vals ...string) (url.Values, []string) {
//}

func CollectTerms(data []map[string]any, facet string) []string {
	var terms [][]string
	for _, item := range data {
		if t, ok := item[facet]; ok {
			terms = append(terms, cast.ToStringSlice(t))
		}
	}
	return lo.Uniq(lo.Flatten(terms))
}

func (f *Facet) AddTerm(term string, ids ...any) *Facet {
	if _, ok := f.Terms[term]; !ok {
		f.Terms[term] = &Term{}
	}
	t := &Term{
		Value: term,
	}
	f.Terms[term] = t
	return f
}

func GetFacetTerms(facet url.Values) []*Term {
	var terms []*Term
	for name, vals := range facet {
		terms = append(terms, NewTerm(name, vals))
	}
	return terms
}

func NewTerm(name string, vals []string) *Term {
	term := &Term{
		Value: name,
		Count: len(vals),
		items: make([]string, len(vals)),
	}
	for i, val := range vals {
		term.items[i] = val
	}
	return term
}

func (t *Term) Bitmap() bitmap.Bitmap {
	var bits bitmap.Bitmap
	for _, item := range t.items {
		bits.Set(cast.ToUint32(item))
	}
	return bits
}
