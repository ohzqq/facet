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
	facets []string
	Facets map[string]url.Values
}

func New(name string, facets []string, data []map[string]any, pk ...string) *Index {
	idx := &Index{
		Name:   name,
		Data:   data,
		PK:     "id",
		facets: facets,
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

func (idx *Index) SetPK(pk string) *Index {
	idx.PK = pk
	return idx
}

func (idx *Index) SetData(data []map[string]any) *Index {
	idx.Data = data
	return idx
}

func (idx *Index) SetFacets(facets []string) *Index {
	idx.facets = facets
	return idx
}
