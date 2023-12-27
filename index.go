package facet

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Index struct {
	Data    []map[string]any  `json:"data,omitempty"`
	Facets  map[string]*Facet `json:"facets"`
	Filters url.Values        `json:"filters"`
}

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

	if idx.Filters != nil {
		return Filter(idx), nil
	}

	return idx, nil
}

func (idx *Index) Filter(q any) *Index {
	filters, err := parseFilters(q)
	if err != nil {
		log.Fatal(err)
	}
	idx.Filters = filters
	return Filter(idx)
}

func (idx *Index) CollectTerms() *Index {
	for name, facet := range idx.Facets {
		facet.Items = make(map[string]*Term)

		vals := collectFacetValues(name, idx.Data)
		for term, ids := range vals {
			facet.Items[term] = NewTerm(term, ids)
		}
	}
	return idx
}

func (idx *Index) GetFacet(name string) (*Facet, error) {
	if f, ok := idx.Facets[name]; ok {
		return f, nil
	}
	return NewFacet(name), errors.New("no such facet")
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

func (idx *Index) String() string {
	return string(idx.JSON())
}

func (idx *Index) JSON() []byte {
	d, err := json.Marshal(idx)
	if err != nil {
		return []byte("{}")
	}
	return d
}

func NewIndexFromReader(r io.Reader) (*Index, error) {
	idx := &Index{}
	err := json.NewDecoder(r).Decode(idx)
	if err != nil {
		return idx, err
	}
	return idx, nil
}

func NewIndexFromFiles(cfg string, data ...string) (*Index, error) {

	idx := &Index{}

	if Exist(cfg) {
		f, err := os.Open(cfg)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		idx, err = NewIndexFromReader(f)
		if err != nil {
			return nil, err
		}
	}

	if len(data) > 0 {
		idx.SetData(lo.ToAnySlice(data)...)
	}
	idx.CollectTerms()
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

func NewIndexFromString(d string) (*Index, error) {
	idx, err := unmarshalIdx([]byte(d))
	if err != nil {
		return nil, err
	}

	if len(idx.Data) > 0 {
		idx.CollectTerms()
	}
	return idx, nil
}

func NewDataFromString(d string) ([]map[string]any, error) {
	return unmarshalData([]byte(d))
}

func NewIndexFromMap(d map[string]any) (*Index, error) {
	idx := &Index{}
	err := mapstructure.Decode(d, idx)
	if err != nil {
		return nil, err
	}
	if len(idx.Data) > 0 {
		idx.CollectTerms()
	}
	return idx, nil
}

func parseCfg(c any) (*Index, error) {
	cfg := &Index{}
	switch val := c.(type) {
	case []byte:
		return unmarshalIdx(val)
	case string:
		if Exist(val) {
			return NewIndexFromFiles(val)
		} else {
			return NewIndexFromString(val)
		}
	case map[string]any:
		return NewIndexFromMap(val)
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

func unmarshalIdx(d []byte) (*Index, error) {
	idx := &Index{}
	err := json.Unmarshal(d, &idx)
	if err != nil {
		return idx, err
	}

	return idx, nil
}

func unmarshalData(d []byte) ([]map[string]any, error) {
	var data []map[string]any
	err := json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
