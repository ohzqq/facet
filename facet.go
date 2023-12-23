package facet

import (
	"net/url"

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

func (f *Facet) AddTerm(term string, ids ...string) *Facet {
	if _, ok := f.terms[term]; !ok {
		f.terms[term] = []string{}
	}
	f.terms[term] = append(f.terms[term], ids...)
	return f
}
