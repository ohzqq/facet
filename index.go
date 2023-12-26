package facet

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Index struct {
	Key    string                `json:"key"`
	Data   []map[string]any      `json:"data,omitempty"`
	items  []any                 `json:"-"`
	facets map[string]url.Values `json:"-"`
	Facets map[string]*Facet     `json:"facets"`
}

type Opt func(*Index) Opt

func New(c any, data ...any) (*Index, error) {
	idx, err := parseCfg(c)
	if err != nil {
		return nil, err
	}

	err = idx.SetData(data...)
	if err != nil {
		return nil, err
	}

	if len(idx.Data) < 1 {
		return idx, errors.New("data is required")
	}

	idx.CollectTerms()

	return idx, nil
}

func (idx *Index) Filter(q url.Values) []map[string]any {
	var bits []*roaring.Bitmap
	for name, filters := range q {
		if facet, ok := idx.Facets[name]; ok {
			bits = append(bits, facet.Filter(filters...))
		}
	}

	filtered := roaring.ParOr(4, bits...)

	ids := filtered.ToArray()
	return FilterItems(idx.Data, lo.ToAnySlice(ids))
}

func (idx *Index) CollectTerms() {
	for name, facet := range idx.Facets {
		facet.Terms = make(map[string]*Term)

		vals := collectFacetValues(name, idx.Data)
		for term, ids := range vals {
			facet.Terms[term] = NewTerm(term, ids)
		}
	}
}

func (idx *Index) GetFacet(name string) (*Facet, error) {
	if f, ok := idx.Facets[name]; ok {
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

func (idx *Index) SetData(data ...any) error {
	for _, datum := range data {
		d, err := parseData(datum)
		if err != nil {
			return err
		}
		idx.Data = append(idx.Data, d...)
	}
	return nil
}

func NewIndexFromFiles(cfg string, data ...string) (*Index, error) {

	idx := &Index{}

	if Exist(cfg) {
		data, err := os.ReadFile(cfg)
		if err != nil {
			return nil, err
		}
		idx, err = unmarshalCfg(data)
		if err != nil {
			return nil, err
		}
	}

	d, err := NewDataFromFiles(data...)
	if err != nil {
		return nil, err
	}
	idx.Data = d
	return idx, nil
}

func NewDataFromFiles(d ...string) ([]map[string]any, error) {
	var data []map[string]any
	for _, datum := range d {
		p, err := dataFromFile(datum)
		if err != nil {
			return nil, err
		}
		data = append(data, p...)
	}
	return data, nil
}

func dataFromFile(d string) ([]map[string]any, error) {
	if Exist(d) {
		data, err := os.ReadFile(d)
		if err != nil {
			return nil, err
		}
		return unmarshalData(data)
	}
	return nil, errors.New("can't read data file")
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func parseCfg(c any) (*Index, error) {
	cfg := &Index{}
	var err error
	switch val := c.(type) {
	case []byte:
		return unmarshalCfg(val)
		if err != nil {
			return nil, err
		}
	case string:
		if Exist(val) {
			return NewIndexFromFiles(val)
		} else {
			return unmarshalCfg([]byte(val))
		}
	case map[string]any:
		if f, ok := val["facets"]; ok {
			cfg.Facets = parseFacetMap(f)
		} else {
			return cfg, errors.New("facets not found in config")
		}
		if data, ok := val["data"]; ok {
			err := cfg.SetData(data)
			if err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}

func unmarshalCfg(d []byte) (*Index, error) {
	cfg := &Index{}
	err := json.Unmarshal(d, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func parseFacetMap(f any) map[string]*Facet {
	facets := make(map[string]*Facet)
	for name, agg := range cast.ToStringMap(f) {
		facet := NewFacet(name)
		err := mapstructure.Decode(agg, facet)
		if err != nil {
			log.Fatal(err)
		}
		facets[name] = facet
	}
	return facets
}

func parseData(d any) ([]map[string]any, error) {
	switch val := d.(type) {
	case []byte:
		return unmarshalData(val)
	case string:
		if Exist(val) {
			return dataFromFile(val)
		} else {
			return unmarshalData([]byte(val))
		}
	case []map[string]any:
		return val, nil
	}
	return nil, errors.New("data couldn't be parsed")
}

func unmarshalData(d []byte) ([]map[string]any, error) {
	var data []map[string]any
	err := json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
