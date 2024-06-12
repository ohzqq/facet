package facet

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

type Facetz struct {
	*Params `json:"params"`
	Facets  []*Field         `json:"facets"`
	data    []map[string]any `json:"hits"`
	ids     []string
	bits    *roaring.Bitmap
}

type Facets struct {
	bits   *roaring.Bitmap
	pk     string
	Fields map[string]*Field
}

func New(data []map[string]any, fields []string, pk string, filters ...any) *Facets {
	f := &Facets{
		bits:   roaring.New(),
		pk:     pk,
		Fields: make(map[string]*Field),
	}
	if len(fields) < 1 && len(data) > 0 {
		fields = maps.Keys(data[0])
	}
	for _, field := range fields {
		f.Fields[field] = NewField(field)
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

func (f Facets) Len() int {
	return int(f.bits.GetCardinality())
}

func (f *Facets) Bitmap() *roaring.Bitmap {
	return f.bits
}
