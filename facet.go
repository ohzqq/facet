package facet

import (
	"encoding/json"
	"net/url"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

type Facets struct {
	fields []*Field
	Attrs  []string `mapstructure:"attributesForFaceting"`
	Data   []map[string]any
	UID    string
}

func New(params any) (*Facets, error) {
	pm := make(map[string]any)
	switch p := params.(type) {
	case []byte:
		err := json.Unmarshal(p, &pm)
		if err != nil {
			return nil, err
		}
	case string:
		q, err := url.ParseQuery(p)
		if err != nil {
			return nil, err
		}
		for attr, vals := range q {
			pm[attr] = vals
		}
	case url.Values:
		for attr, vals := range p {
			pm[attr] = vals
		}
	case map[string]any:
		pm = p
	}

	facets := &Facets{}
	err := mapstructure.Decode(pm, facets)
	if err != nil {
		return nil, err
	}

	return facets, nil
}

func NewFacets(data []map[string]any, attrs []string) *Facets {
	return &Facets{
		UID:   "id",
		Attrs: attrs,
		Data:  data,
	}
}

func (f *Facets) Calculate() *Facets {
	facets := CalculateFacets(f.Data, f.Attrs, f.UID)
	f.fields = facets
	return f
}

func CalculateFacets(data []map[string]any, fields []string, ident ...string) []*Field {
	facets := NewFields(fields)

	uid := "id"
	if len(ident) > 0 {
		uid = ident[0]
	}

	for id, d := range data {
		if i, ok := d[uid]; ok {
			id = cast.ToInt(i)
		}
		for _, facet := range facets {
			if val, ok := d[facet.Attribute]; ok {
				facet.Add(
					val,
					[]int{id},
				)
			}
		}
	}
	return facets
}
