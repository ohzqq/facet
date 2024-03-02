package facet

import (
	"encoding/json"
	"log"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Facets struct {
	*Params `json:"params"`
	Facets  []*Field         `json:"facets"`
	Hits    []map[string]any `json:"hits"`
	bits    *roaring.Bitmap
	hits    []int
}

func New(params any) (*Facets, error) {
	var err error

	p, err := ParseParams(params)
	if err != nil {
		return nil, err
	}

	facets := NewFacets(p.Attrs())
	facets.Params = p

	facets.Hits, err = facets.Data()
	if err != nil {
		log.Fatal(err)
	}

	facets.Calculate()

	if facets.vals.Has("facetFilters") {
		filtered, err := facets.Filter(facets.Filters())
		if err != nil {
			return nil, err
		}
		return filtered.Calculate(), nil
	}

	return facets, nil
}

func NewFacets(fields []string) *Facets {
	return &Facets{
		bits:   roaring.New(),
		Facets: NewFields(fields),
	}
}

func (f *Facets) Calculate() *Facets {
	for id, d := range f.Hits {
		f.bits.AddInt(id)
		for _, facet := range f.Facets {
			if val, ok := d[facet.Attribute]; ok {
				facet.Add(
					val,
					[]int{id},
				)
			}
		}
	}
	return f
}

func (f *Facets) Filter(filters []any) (*Facets, error) {
	filtered, err := Filter(f.bits, f.Facets, filters)
	if err != nil {
		return nil, err
	}

	facets := NewFacets(f.Attrs())
	facets.Params = f.Params

	facets.hits = f.FilterBits(filtered)

	if len(facets.hits) > 0 {
		facets.Hits = f.GetByID(facets.hits...)
	}

	return facets, nil
}

func (f *Facets) FilterBits(bits *roaring.Bitmap) []int {
	bits.And(f.bits)
	return cast.ToIntSlice(bits.ToArray())
}

func (f *Facets) GetByID(ids ...int) []map[string]any {
	var res []map[string]any
	for id, d := range f.Hits {
		if f.bits.ContainsInt(id) {
			res = append(res, d)
		}
	}
	return res
}

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
}

func (f Facets) Count() int {
	return len(f.hits)
}

func (f Facets) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facets) MarshalJSON() ([]byte, error) {
	facets := make(map[string]*Field)
	for _, facet := range f.Facets {
		facets[facet.Attribute] = facet
	}

	enc := make(map[string]any)
	enc["params"] = f.EncodeQuery()
	enc["facets"] = facets
	enc["hits"] = f.Hits

	return json.Marshal(enc)
}
