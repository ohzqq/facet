package facet

import (
	"encoding/json"
	"log"

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

	facets := NewFacets(f.Attrs())
	facets.Params = f.Params

	facets.Hits = f.FilterBits(filtered)

	if len(facets.Hits) > 0 {
		facets.data = f.GetByID(facets.Hits...)
	}

	return facets, nil
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
