package facet

import (
	"encoding/json"

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

	facets.data, err = facets.Data()
	if err != nil {
		return nil, err
	}

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

//func (f *Facets) GetByID(ids ...string) []map[string]any {
//}

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
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
