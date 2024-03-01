package facet

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
)

type Params struct {
	FacetFilters []string `mapstructure:"facetFilters"`
	DataSrc      []string `mapstructure:"data"`
	FacetAttrs   []string `mapstructure:"attributesForFaceting"`
	Ident        []string `mapstructure:"uid"`
	vals         url.Values
}

func ParseParams(params any) (*Params, error) {
	p := &Params{}

	var err error
	switch param := params.(type) {
	case string:
		p.vals, err = url.ParseQuery(param)
		if err != nil {
			return nil, err
		}
	case url.Values:
		p.vals = param
	}

	return p, nil
}

func (f Params) UID() string {
	if f.vals.Has("uid") {
		return f.vals.Get("uid")
	}
	return "id"
}

func (f Params) Attrs() []string {
	if f.vals.Has("attributesForFaceting") {
		attrs := f.vals["attributesForFaceting"]
		if len(attrs) == 1 {
			return strings.Split(f.vals.Get("attributesForFaceting"), ",")
		}
		return attrs
	}
	return []string{}
}

func (p Params) Data() ([]map[string]any, error) {
	var data []map[string]any

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

func DecodeData(r io.Reader, data *[]map[string]any) error {
	dec := json.NewDecoder(r)
	for {
		m := make(map[string]any)
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*data = append(*data, m)
	}
	return nil
}
