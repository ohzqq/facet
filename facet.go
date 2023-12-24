package facet

import (
	"net/url"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Facet struct {
	Name  string
	terms url.Values
}

func NewFacet(name string, pk string, data []map[string]any) url.Values {
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

func CollectTerms(facet string, data []map[string]any) []string {
	var terms [][]string
	for _, item := range data {
		if t, ok := item[facet]; ok {
			terms = append(terms, cast.ToStringSlice(t))
		}
	}
	return lo.Uniq(lo.Flatten(terms))
}

func (f *Facet) AddTerm(term string, ids ...string) *Facet {
	if _, ok := f.terms[term]; !ok {
		f.terms[term] = []string{}
	}
	f.terms[term] = append(f.terms[term], ids...)
	return f
}
