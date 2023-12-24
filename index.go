package facet

import (
	"net/url"

	"github.com/kelindar/bitmap"
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
		idx.Facets[f] = NewFacetVals(f, idx.PK, data)
	}
	return idx
}

func (idx *Index) Bitmap(ids ...any) bitmap.Bitmap {
	if len(ids) > 0 {
		return NewBitmap(ids)
	}
	return NewBitmap(lo.ToAnySlice(idx.items))
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

func Filter(idx *Index, facet string, op string, filters []string, ids ...any) (*Index, string, string, []string) {
	//agg := idx.GetFacet(facet)
	bitIDs := idx.Bitmap(ids...)
	f := lo.Slice(filters, 0, 1)
	if len(f) > 0 {
		term := idx.GetTerm(facet, f[0])
		switch op {
		case "and":
			bitIDs.And(term.Bitmap())
		case "or":
			bitIDs.Or(term.Bitmap())
		}
		var rest []map[string]any
		for _, item := range idx.Data {
			if bitIDs.Contains(cast.ToUint32(item[idx.PK])) {
				rest = append(rest, item)
			}
		}
		idx.Data = rest
		return Filter(idx, facet, op, filters[0:])
	}

	return idx, facet, op, filters
}

func CollectIDsInt(pk string, data []map[string]any) []uint32 {
	iter := func(item map[string]any, _ int) uint32 {
		return cast.ToUint32(item[pk])
	}
	return lo.Map(data, iter)
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

func (idx *Index) GetTerm(facet, term string) *Term {
	for _, t := range idx.GetTerms(facet) {
		if t.Value == term {
			return t
		}
	}
	return &Term{Value: term}
}

func (idx *Index) GetTerms(name string) []*Term {
	return GetFacetTerms(idx.GetFacet(name))
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
