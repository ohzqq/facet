package facet

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
)

type Facets struct {
	Facets []*Field
	data   []map[string]any
	*Params
}

func New(params any) (*Facets, error) {
	facets := NewFacets()

	var err error

	facets.Params, err = ParseParams(params)
	if err != nil {
		return nil, err
	}

	facets.data, err = facets.Data()
	if err != nil {
		return nil, err
	}

	return facets, nil
}

func NewFacets() *Facets {
	return &Facets{}
}

func (f Facets) GetFacet(attr string) *Field {
	for _, facet := range f.Facets {
		if facet.Attribute == attr {
			return facet
		}
	}
	return &Field{}
}

func (f Facets) EncodeQuery() string {
	return f.vals.Encode()
}

func (f *Facets) Calculate() *Facets {
	facets := CalculateFacets(f.data, f.Attrs(), f.UID())
	f.Facets = facets
	return f
}

func (f *Facets) MarshalJSON() ([]byte, error) {
	facets := make(map[string]any)
	facets["params"] = f.EncodeQuery()
	facets["facets"] = f.Facets

	return json.Marshal(facets)
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
				if auth := cast.ToString(val); auth == "mafia" {
					fmt.Printf("%+v\n", val)
				}
				facet.Add(
					val,
					[]int{id},
				)
			}
		}
	}
	return facets
}
