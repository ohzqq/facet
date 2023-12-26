package facet

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"github.com/RoaringBitmap/roaring"
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

type Opt func(*Index) Opt

func New(cfg, data any) (*Index, error) {
	idx := &Index{}
	return idx, nil
}

func NewFromFiles(c string, data ...string) (*Index, error) {
	idx := &Index{}
	key, facets, err := parseCfg(c)
	if err != nil {
		return idx, err
	}
	idx.Key = key
	idx.FacetCfg = facets

	for _, datum := range data {
		d, err := os.ReadFile(datum)
		if err != nil {
			return idx, err
		}
		var dd []map[string]any
		err = json.Unmarshal(d, &dd)
		if err != nil {
			return idx, err
		}
		idx.Data = append(idx.Data, dd...)
	}

	if len(idx.Data) < 1 {
		return idx, errors.New("data is required")
	}

	idx.Facets()

	return idx, nil
}

func NewIdx(name string, facets []string, data []map[string]any, pk ...string) *Index {
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

func (idx *Index) Roar() *roaring.Bitmap {
	ids := CollectIDsInt(idx.Key, idx.Data)
	return roaring.BitmapOf(ids...)
}

func (idx *Index) Filter(q url.Values) []map[string]any {
	var bits []*roaring.Bitmap
	for name, filters := range q {
		if facet, ok := idx.FacetCfg[name]; ok {
			bits = append(bits, facet.Filter(filters...))
		}
	}

	filtered := roaring.ParOr(4, bits...)

	ids := filtered.ToArray()
	items := make([]map[string]any, len(ids))
	for _, item := range idx.Data {
		for i, id := range ids {
			if cast.ToString(id) == cast.ToString(item[idx.Key]) {
				items[i] = item
			}
		}
	}
	return items
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

func (idx *Index) GetTerm(facet, term string) (*Term, error) {
	f, err := idx.GetFacet(facet)
	if err != nil {
		return nil, err
	}

	return f.GetTerm(term), nil
}

func (idx *Index) SetPK(pk string) *Index {
	idx.Key = pk
	return idx
}

func (idx *Index) SetData(data []map[string]any) *Index {
	idx.Data = data
	return idx
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
