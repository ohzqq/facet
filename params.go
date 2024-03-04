package facet

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
)

type Params struct {
	vals url.Values
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

func (p Params) Attrs() []string {
	if p.vals.Has("attributesForFaceting") {
		attrs := p.vals["attributesForFaceting"]
		if len(attrs) == 1 {
			return strings.Split(p.vals.Get("attributesForFaceting"), ",")
		}
		return attrs
	}
	return []string{}
}

func (f Params) UID() string {
	if f.vals.Has("uid") {
		return f.vals.Get("uid")
	}
	return ""
}

func (p *Params) Filters() []any {
	if p.vals.Has("facetFilters") {
		fils, err := unmarshalFilter(p.vals.Get("facetFilters"))
		if err != nil {
		}
		return fils
	}
	return []any{}
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
	//fmt.Printf("num data %d\n", len(data))

	return data, nil
}

func (p Params) MarshalJSON() ([]byte, error) {
	params := p.vals.Encode()
	return json.Marshal(params)
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
