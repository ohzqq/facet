package facet

import (
	"encoding/json"
	"log"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Facets struct {
	Facets []*Field
	bits   *roaring.Bitmap
	Hits   []int
	data   []map[string]any
	*Params
}

func New(params any) (*Facets, error) {
	facets := NewFacets()

	var err error

	facets.Params, err = ParseParams(params)
	if err != nil {
		return nil, err
	}

	facets.Facets = NewFields(facets.Attrs())

	facets.Calculate()

	if facets.vals.Has("facetFilters") {
		filters := facets.Filters()
		facets, err = facets.Filter(filters)
		if err != nil {
			return nil, err
		}
		facets.Calculate()
	}

	return facets, nil
}

func NewFacets() *Facets {
	return &Facets{
		bits: roaring.New(),
	}
}

func (f *Facets) Calculate() *Facets {
	var err error
	f.data, err = f.Data()
	if err != nil {
		log.Fatal(err)
	}

	uid := f.UID()

	for id, d := range f.data {
		if i, ok := d[uid]; ok {
			id = cast.ToInt(i)
		}
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

	f.Hits = f.FilterBits(filtered)

	return f, nil
}

func (f *Facets) FilterBits(bits *roaring.Bitmap) []int {
	bits.And(f.bits)
	return cast.ToIntSlice(bits.ToArray())
}

func (f *Facets) GetByID(ids ...int) []map[string]any {
	uid := f.UID()
	var res []map[string]any
	for id, d := range f.data {
		if i, ok := d[uid]; ok {
			id = cast.ToInt(i)
		}
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

func (p Facets) Data() ([]map[string]any, error) {
	var data []map[string]any

	if len(p.Hits) > 0 {
		return p.GetByID(p.Hits...), nil
	}

	if p.vals.Has("data") {
		for _, file := range p.vals["data"] {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			err = DecodeData(f, &data)
			if err != nil {
				return nil, err
			}
		}
	}

	return data, nil
}

func (f Facets) Count() int {
	return len(f.Hits)
}

func (f Facets) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facets) MarshalJSON() ([]byte, error) {
	facets := make(map[string]any)
	facets["params"] = f.EncodeQuery()
	facets["facets"] = f.Facets

	return json.Marshal(facets)
}
