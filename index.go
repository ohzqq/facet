package facet

import (
	"net/url"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

type Index struct {
	Name   string
	PK     string
	Data   []map[string]any
	items  []string
	Facets map[string]url.Values
}

func New(name string, facets []string, data []map[string]any, pk ...string) *Index {
	idx := &Index{
		Name:   name,
		Data:   data,
		PK:     "id",
		Facets: make(map[string]url.Values),
	}
	if len(pk) > 0 {
		idx.PK = pk[0]
	}
	idx.items = CollectIDs(idx.PK, data)
	for _, f := range facets {
		idx.Facets[f] = NewFacet(f, idx.PK, data)
	}
	return idx
}

func (idx *Index) GetByID(ids []string) []map[string]any {
	var data []map[string]any
	for _, item := range idx.Data {
		if lo.Contains(ids, cast.ToString(item[idx.PK])) {
			data = append(data, item)
		}
	}
	return data
}

func CollectIDs(pk string, data []map[string]any) []string {
	iter := func(item map[string]any, _ int) string {
		return cast.ToString(item[pk])
	}
	return lo.Map(data, iter)
}

func (idx *Index) GetFacet(name string) url.Values {
	if f, ok := idx.Facets[name]; ok {
		return f
	}
	return url.Values{}
}

func (idx *Index) GetFacetValues(name string) []string {
	return maps.Keys(idx.GetFacet(name))
}

func (idx *Index) GetFacetTermItems(facet, term string) []string {
	f := idx.GetFacet(facet)
	if f.Has(term) {
		return f[term]
	}
	return []string{}
}

func (idx *Index) SetPK(pk string) *Index {
	idx.PK = pk
	return idx
}

func (idx *Index) SetData(data []map[string]any) *Index {
	idx.Data = data
	return idx
}
