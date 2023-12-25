package facet

import (
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/kelindar/bitmap"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Index struct {
	Name     string                `json:"name"`
	Key      string                `json:"key"`
	Data     []map[string]any      `json:"-"`
	items    []string              `json:"-"`
	facets   map[string]url.Values `json:"-"`
	FacetCfg map[string]*Facet     `json:"facets"`
}

func New(name string, facets []string, data []map[string]any, pk ...string) *Index {
	idx := &Index{
		Name:   name,
		Data:   data,
		Key:    "id",
		facets: make(map[string]url.Values),
	}
	if len(pk) > 0 {
		idx.Key = pk[0]
	}
	idx.items = CollectIDs(idx.Key, data)
	for _, f := range facets {
		idx.facets[f] = CollectFacetValues(f, idx.Key, data)
	}
	return idx
}

func (idx *Index) Facets() map[string]*Facet {
	//ids := CollectIDsInt(idx.Key, idx.Data)
	//items := NewBitmap(lo.ToAnySlice(idx.items))
	//facets := make(map[string]url.Values)
	idx.CollectTerms()
	return idx.FacetCfg
}

func (idx *Index) Filter(q url.Values) []string {
	og := CollectAnyIDs(idx.Key, idx.Data)
	items := NewBitmap(og)

	var op string
	var bits []bitmap.Bitmap
	for name, filters := range q {
		if facet, ok := idx.FacetCfg[name]; ok {
			for _, filter := range filters {
				term := idx.GetTerm(name, filter)
				b := term.Bitmap()
				bits = append(bits, b)
				switch facet.Operator {
				case "or":
					op = "or"
				case "and":
					op = "and"
				}
			}
		}
	}

	switch len(bits) {
	case 0:
	case 1:
		switch op {
		case "or":
			items.Or(bits[0])
		case "and":
			items.And(bits[0])
		}
	default:
		switch op {
		case "or":
			items.Or(bits[0], bits[1:]...)
		case "and":
			items.And(bits[0], bits[1:]...)
		}
	}

	var ids []string
	items.Range(func(x uint32) {
		id := strconv.Itoa(int(x))
		ids = append(ids, id)
	})
	return ids
}

func or(bits ...bitmap.Bitmap) []bitmap.Bitmap {
	if len(bits) < 2 {
		return bits
	}
	bits[0].Or(bits[1])
	return bits[2:]
}

func and(bits ...bitmap.Bitmap) []bitmap.Bitmap {
	if len(bits) < 2 {
		return bits
	}
	bits[0].And(bits[1])
	return bits[2:]
}

func (idx *Index) CollectTerms() {
	for name, facet := range idx.FacetCfg {
		facet.Terms = make(map[string]*Term)

		vals := CollectFacetValues(name, idx.Key, idx.Data)
		for term, ids := range vals {
			facet.Terms[term] = NewTerm(term, ids)
		}
	}
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
		if lo.Contains(ids, cast.ToString(item[idx.Key])) {
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
			if bitIDs.Contains(cast.ToUint32(item[idx.Key])) {
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

func CollectAnyIDs(pk string, data []map[string]any) []any {
	iter := func(item map[string]any, _ int) any {
		return item[pk]
	}
	return lo.Map(data, iter)
}

func CollectIDs(pk string, data []map[string]any) []string {
	iter := func(item map[string]any, _ int) string {
		return cast.ToString(item[pk])
	}
	return lo.Map(data, iter)
}

func (idx *Index) GetFacet(name string) (*Facet, error) {
	if f, ok := idx.FacetCfg[name]; ok {
		return f, nil
	}
	return &Facet{}, errors.New("no such facet")
}

func (idx *Index) GetTerm(facet, term string) *Term {
	f, err := idx.GetFacet(facet)
	if err != nil {
		log.Fatal(err)
	}

	return f.GetTerm(term)
}

func (idx *Index) SetPK(pk string) *Index {
	idx.Key = pk
	return idx
}

func (idx *Index) SetData(data []map[string]any) *Index {
	idx.Data = data
	return idx
}
