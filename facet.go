package facet

import (
	"encoding/json"
	"log"

	"github.com/RoaringBitmap/roaring"
)

type Facets struct {
	*Params `json:"params"`
	Facets  []*Field         `json:"facets"`
	Hits    []map[string]any `json:"hits"`
	bits    *roaring.Bitmap
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

	f.bits.And(filtered)

	if f.bits.GetCardinality() > 0 {
		facets.Hits = ItemsByBitmap(f.Hits, f.bits)
	}

	return facets, nil
}

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
}

func (f Facets) Len() int {
	return int(f.bits.GetCardinality())
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

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}
