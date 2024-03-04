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
	data    []map[string]any `json:"hits"`
	ids     []string
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

	facets.data, err = facets.Data()
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
	for id, d := range f.data {
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
		facets.data = f.filteredItems()
	}

	return facets, nil
}

func (f *Facets) filteredItems() []map[string]any {
	var res []map[string]any
	f.bits.Iterate(func(x uint32) bool {
		res = append(res, f.data[int(x)])
		return true
	})
	return res
}

func (f *Facets) getHits() []any {
	var res []any
	var uid string
	if f.vals.Has("uid") {
		uid = f.vals.Get("uid")
	}

	f.bits.Iterate(func(x uint32) bool {
		d := f.data[int(x)]
		if id, ok := d[uid]; ok {
			res = append(res, id)
		} else {
			res = append(res, int(x))
		}
		return true
	})
	return res
}

func (f Facets) getItem(id int) (map[string]any, bool) {
	uid := f.UID()
	for idx, d := range f.data {
		if i, ok := d[uid]; ok {
			if cast.ToInt(i) == id {
				return d, true
			}
		} else if id == idx {
			return d, true
		}
	}
	return nil, false
}

func (f Facets) getItems() []map[string]any {
	//uid := f.UID()

	ids := cast.ToIntSlice(f.bits.ToArray())

	data := make([]map[string]any, len(ids))
	for _, id := range ids {
		if d, ok := f.getItem(id); ok {
			data = append(data, d)
		}
	}

	return data
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

func (f *Facets) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Facets) MarshalJSON() ([]byte, error) {
	facets := make(map[string]*Field)
	for _, facet := range f.Facets {
		facets[facet.Attribute] = facet
	}

	enc := make(map[string]any)
	enc["params"] = f.EncodeQuery()
	enc["facets"] = facets
	enc["hits"] = f.getHits()
	enc["nbHits"] = f.Len()

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
