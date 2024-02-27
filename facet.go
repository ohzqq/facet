package srch

import "github.com/ohzqq/srch/txt"

type Facet struct {
	*Field
}

func NewFacet(attr string) *Facet {
	f := NewField(attr)
	f.SetAnalyzer(txt.Keyword())

	return &Facet{
		Field: f,
	}
}

func NewFacets(attrs []string) map[string]*Facet {
	fields := make(map[string]*Facet)
	for _, attr := range attrs {
		fields[attr] = NewFacet(attr)
	}

	return fields
}

func CalculateFacets(data []map[string]any, fields []string) map[string]*Facet {
	facets := NewFacets(fields)
	for id, d := range data {
		for attr, facet := range facets {
			if val, ok := d[attr]; ok {
				facet.Add(val, []int{id})
			}
		}
	}
	return facets
}
