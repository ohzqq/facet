package facet

import (
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Index struct {
	Name   string
	PK     string
	Data   []map[string]any
	items  []string
	Facets []string
	facets map[string]*Facet
}

func New(name string, facets []string, data ...map[string]any) *Index {
	idx := &Index{
		Name:   name,
		Data:   data,
		PK:     "id",
		Facets: facets,
		items:  make([]string, len(data)),
		facets: make(map[string]*Facet),
	}
	for _, f := range facets {
		idx.facets[f] = NewFacet(f)
	}
	return idx
}

func (idx *Index) getIDs() []string {
	iter := func(item map[string]any, _ int) string {
		return cast.ToString(item[idx.PK])
	}
	return lo.Map(idx.Data, iter)
}

func getTerms(data []map[string]any, term string) map[string]any {
	trans := func(item map[string]any) (string, any) {
		if t, ok := item[term]; ok {
		}
		return "", nil
	}
}

func (idx *Index) processData() *Index {
	for i, item := range idx.Data {
		idx.items[i] = cast.ToString(item[idx.PK])
		for _, f := range idx.Facets {
			facet := idx.GetFacet(f)
			if terms, ok := item[f]; ok {
				for _, term := range terms.([]any) {
					t := cast.ToString(term)
					facet.AddTerm(t, cast.ToString(item[idx.PK]))
				}
			}
		}
	}
	return idx
}

func (idx *Index) GetFacet(name string) *Facet {
	if f, ok := idx.facets[name]; ok {
		return f
	}
	idx.facets[name] = NewFacet(name)
	return idx.facets[name]
}

//func (idx *Index) ProcessFacets() *Index {
//  for _, f := range idx.Facets {
//    var vals [][]any
//    for _, b := range idx.Data {
//      vals = append(vals, b[f].([]any))
//    }
//    idx.facets[f] = lo.Uniq(lo.Flatten(vals))
//  }
//  return idx
//}

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
