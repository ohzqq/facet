package facet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("workers", 1)
}

// Index is a structure for facets and data.
type Index struct {
	Data    []map[string]any `json:"data,omitempty"`
	Facets  []*Facet         `json:"facets"`
	Filters url.Values       `json:"filters"`
}

// New initializes an index.
func New(c any, data ...any) (*Index, error) {
	idx, err := parseCfg(c)
	if err != nil {
		return nil, err
	}

	err = idx.SetData(data...)
	if err != nil {
		return nil, err
	}

	if idx.Filters != nil {
		return Filter(idx), nil
	}

	return idx, nil
}

// Filter idx.Data and re-calculate facets.
func (idx *Index) Filter(q any) *Index {
	filters, err := ParseFilters(q)
	if err != nil {
		log.Fatal(err)
	}

	idx.Filters = filters
	return Filter(idx)
}

// CollectItems collects a facet's items from the data set.
func (idx *Index) CollectItems() *Index {
	for _, facet := range idx.Facets {
		facet.CollectItems(idx.Data)
	}
	return idx
}

// GetConfig returns a map of the Index's config.
func (idx *Index) GetConfig() map[string]any {
	var facets []map[string]any
	for _, f := range idx.Facets {
		facets = append(facets, f.GetConfig())
	}
	return map[string]any{
		"facets": facets,
	}
}

// GetFacet returns a facet.
func (idx *Index) GetFacet(name string) *Facet {
	for _, facet := range idx.Facets {
		if facet.Attribute == name {
			return facet
		}
	}
	return NewFacet(name)
}

// SetData sets the data set for the index.
func (idx *Index) SetData(data ...any) error {
	for _, datum := range data {
		d, err := parseData(datum)
		if err != nil {
			return err
		}
		idx.Data = append(idx.Data, d...)
	}
	idx.CollectItems()
	return nil
}

// Decode unmarshals json from an io.Reader.
func (idx *Index) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(idx)
	if err != nil {
		return err
	}
	return nil
}

// Encode marshals json from an io.Writer.
func (idx *Index) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(idx)
}

// DecodeData unmarshals data from an io.Reader.
func (idx *Index) DecodeData(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&idx.Data)
	if err != nil {
		return err
	}
	return nil
}

// String returns an Index as a json formatted string.
func (idx *Index) String() string {
	return string(idx.JSON())
}

// JSON marshals an Index to json.
func (idx *Index) JSON() []byte {
	var buf bytes.Buffer
	err := idx.Encode(&buf)
	if err != nil {
		return []byte("{}")
	}
	return buf.Bytes()
}

// Print writes Index json to stdout.
func (idx *Index) Print() {
	err := idx.Encode(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// PrettyPrint writes Index indented json to stdout.
func (idx *Index) PrettyPrint() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(idx)
	if err != nil {
		log.Fatal(err)
	}
}

// DecodeData decodes data from a io.Reader.
func DecodeData(r io.Reader) ([]map[string]any, error) {
	var data []map[string]any
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// NewIndexFromFiles initializes an index from files.
func NewIndexFromFiles(cfg string) (*Index, error) {
	idx := &Index{}

	f, err := os.Open(cfg)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = idx.Decode(f)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

// NewDataFromFiles parses index data from files.
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
	data, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return DecodeData(data)
}

// NewIndexFromString initializes an index from a json formatted string.
func NewIndexFromString(d string) (*Index, error) {
	idx := &Index{}
	buf := bytes.NewBufferString(d)
	err := idx.Decode(buf)
	if err != nil {
		return nil, err
	}

	if len(idx.Data) > 0 {
		idx.CollectItems()
	}
	return idx, nil
}

// NewDatFromString parses index data from a json formatted string.
func NewDataFromString(d string) ([]map[string]any, error) {
	buf := bytes.NewBufferString(d)
	return DecodeData(buf)
}

// NewIndexFromMap initalizes an index from a map[string]any.
func NewIndexFromMap(d map[string]any) (*Index, error) {
	idx := &Index{}
	err := mapstructure.Decode(d, idx)
	if err != nil {
		return nil, err
	}
	if len(idx.Data) > 0 {
		idx.CollectItems()
	}
	return idx, nil
}

func parseCfg(c any) (*Index, error) {
	cfg := &Index{}
	switch val := c.(type) {
	case []byte:
		buf := bytes.NewBuffer(val)
		err := cfg.Decode(buf)
		return cfg, err
	case string:
		if exist(val) {
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
		if exist(val) {
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

func exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
