package facet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cast"
)

type Facets struct {
	Facets []*Field
	data   []map[string]any
	Params url.Values
}

func New(params any) (*Facets, error) {
	facets := NewFacets()

	var err error
	switch p := params.(type) {
	case string:
		facets.Params, err = url.ParseQuery(p)
		if err != nil {
			return nil, err
		}
	case url.Values:
		facets.Params = p
	}

	if facets.Params.Has("data") {
		for _, file := range facets.Params["data"] {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			err = facets.DecodeData(f)
			if err != nil {
				return nil, err
			}
		}
	}

	return facets, nil
}

func NewFacets() *Facets {
	return &Facets{}
}

func (f Facets) UID() string {
	if f.Params.Has("uid") {
		return f.Params.Get("uid")
	}
	return "id"
}

func (f Facets) Attrs() []string {
	if f.Params.Has("attributesForFaceting") {
		attrs := f.Params["attributesForFaceting"]
		if len(attrs) == 1 {
			return strings.Split(attrs[0], ",")
		}
		return attrs
	}
	return []string{}
}

func (f *Facets) DecodeData(r io.Reader) error {
	dec := json.NewDecoder(r)
	for {
		m := make(map[string]any)
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		f.data = append(f.data, m)
	}
	return nil
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
	return f.Params.Encode()
}

func (f *Facets) Calculate() *Facets {
	facets := CalculateFacets(f.data, f.Attrs(), f.UID())
	f.Facets = facets
	return f
}

func (f *Facets) Search(q any) *Facets {
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
