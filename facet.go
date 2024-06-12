package facet

import (
	"encoding/json"
	"io"
	"log"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

type Facetz struct {
	*Params `json:"params"`
	Facets  []*Fieldz        `json:"facets"`
	data    []map[string]any `json:"hits"`
	ids     []string
	bits    *roaring.Bitmap
}

type Facets struct {
	bits   *roaring.Bitmap
	pk     string
	Fields map[string]*Fieldz
}

func New(data []map[string]any, fields []string, pk string, filters ...any) *Facets {
	f := &Facets{
		bits:   roaring.New(),
		pk:     pk,
		Fields: make(map[string]*Fieldz),
	}
	if len(fields) < 1 && len(data) > 0 {
		fields = maps.Keys(data[0])
	}
	for _, field := range fields {
		f.Fields[field] = NewFieldz(field)
	}

	for idx, d := range data {
		if id, ok := d[pk]; ok {
			idx = cast.ToInt(id)
		}
		f.bits.AddInt(idx)

		facetable := lo.PickByKeys(d, fields)
		if len(facetable) < 1 {
			facetable = d
		}
		for field, v := range facetable {
			f.Fields[field].Add(v, []int{idx})
		}
	}

	if len(filters) > 0 {
		err := f.Filter(filters)
		if err != nil {
			return nil
		}
		if f.bits.GetCardinality() > 0 {
			res := lo.Filter(data, f.filterItem)
			return New(res, fields, pk)
		}
	}
	return f
}

func (f *Facets) filterItem(item map[string]any, idx int) bool {
	if id, ok := item[f.pk]; ok {
		idx = cast.ToInt(id)
	}
	return f.bits.ContainsInt(idx)
}

func (f *Facets) Filter(filters []any) error {
	fields := maps.Values(f.Fields)
	filtered, err := Filter(f.bits, fields, filters)
	if err != nil {
		return err
	}

	f.bits.And(filtered)

	return nil
}

func NNew(params any) (*Facetz, error) {
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

func NewFacets(fields []string) *Facetz {
	return &Facetz{
		bits:   roaring.New(),
		Facets: NewFieldzz(fields),
	}
}

func (f *Facetz) Calculate() *Facetz {
	var uid string
	if f.vals.Has("uid") {
		uid = f.vals.Get("uid")
	}

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

func (f *Facetz) Filter(filters []any) (*Facetz, error) {
	filtered, err := Filter(f.bits, f.Facets, filters)
	if err != nil {
		return nil, err
	}

	facets := NewFacets(f.Attrs())
	facets.Params = f.Params

	f.bits.And(filtered)

	if f.bits.GetCardinality() > 0 {
		facets.data = f.getHits()
	}

	return facets, nil
}

func (f Facetz) getHits() []map[string]any {
	var uid string
	if f.vals.Has("uid") {
		uid = f.vals.Get("uid")
	}
	var hits []map[string]any
	for idx, d := range f.data {
		if i, ok := d[uid]; ok {
			idx = cast.ToInt(i)
		}
		if f.bits.ContainsInt(idx) {
			hits = append(hits, d)
		}
	}
	return hits
}

func (f Facetz) GetFacet(attr string) *Fieldz {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Fieldz{}
}

func (f Facetz) Len() int {
	return int(f.bits.GetCardinality())
}

func (f Facetz) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facetz) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Facetz) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	for _, d := range f.data {
		err := enc.Encode(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Facetz) MarshalJSON() ([]byte, error) {
	enc := f.resultMeta()
	enc["hits"] = f.data
	if len(f.data) < 1 {
		enc["hits"] = []any{}
	}

	return json.Marshal(enc)
}

func (f *Facetz) resultMeta() map[string]any {
	enc := make(map[string]any)

	facets := make(map[string]*Fieldz)
	for _, facet := range f.Facets {
		facets[facet.Attribute] = facet
	}
	enc["facets"] = facets

	f.vals.Set("nbHits", cast.ToString(f.Len()))
	enc["params"] = f.EncodeQuery()
	return enc
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}
