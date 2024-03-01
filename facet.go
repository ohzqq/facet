package facet

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Facets struct {
	Facets []*Field
	bits   *roaring.Bitmap
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

	return facets, nil
}

func NewFacets() *Facets {
	return &Facets{
		bits: roaring.New(),
	}
}

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
}

func (f Facets) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facets) Calculate() *Facets {
	uid := f.UID()

	for id, d := range f.data {
		if i, ok := d[uid]; ok {
			id = cast.ToInt(i)
		}
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

func (f *Facets) MarshalJSON() ([]byte, error) {
	facets := make(map[string]any)
	facets["params"] = f.EncodeQuery()
	facets["facets"] = f.Facets

	return json.Marshal(facets)
}

func (idx Facets) Bitmap() *roaring.Bitmap {
	bits := roaring.New()
	bits.AddRange(0, uint64(len(idx.data)))
	return bits
}
