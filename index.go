package facet

import "github.com/samber/lo"

type Index struct {
	Name   string
	PK     string
	Data   []map[string]any
	Facets []string
	facets map[string][]any
}

func New(name string, facets []string, data ...map[string]any) *Index {
	return &Index{
		Name:   name,
		Data:   data,
		PK:     "id",
		Facets: facets,
		facets: make(map[string][]any),
	}
}

func (idx *Index) SetPK(pk string) *Index {
	idx.PK = pk
	return idx
}

func (idx *Index) SetData(data []map[string]any) *Index {
	idx.Data = data
	return idx
}

func (idx *Index) SetFacets(facets []string) *Index {
	idx.Facets = facets
	return idx
}

func (idx *Index) ProcessFacets() *Index {
	for _, f := range idx.Facets {
		var vals [][]any
		for _, b := range idx.Data {
			vals = append(vals, b[f].([]any))
		}
		idx.facets[f] = lo.Uniq(lo.Flatten(vals))
	}
	return idx
}
