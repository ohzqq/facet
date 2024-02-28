package facet

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cast"
)

type Facets struct {
	Facets []*Field `json:"facets"`
	data   []map[string]any
	Query  url.Values
}

func New(params any) (*Facets, error) {
	facets := NewFacets()

	var err error
	switch p := params.(type) {
	case string:
		facets.Query, err = url.ParseQuery(p)
		if err != nil {
			return nil, err
		}
		//valsToMap(pm, facets.Query)
	case url.Values:
		facets.Query = p
		//valsToMap(pm, facets.Query)
	}

	if facets.Query.Has("data") {
		for _, file := range facets.Query["data"] {
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
	if f.Query.Has("uid") {
		return f.Query.Get("uid")
	}
	return "id"
}

func (f Facets) Attrs() []string {
	if f.Query.Has("attributesForFaceting") {
		attrs := f.Query["attributesForFaceting"]
		if len(attrs) == 1 {
			return strings.Split(attrs[0], ",")
		}
		return attrs
	}
	return []string{}
}

func (f *Facets) Calculate() *Facets {
	facets := CalculateFacets(f.data, f.Attrs(), f.UID())
	f.Facets = facets
	return f
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

func (f Facets) EncodeQuery() string {
	if f.Query == nil {
		f.Query = make(url.Values)
		f.Query.Set("uid", f.UID())
		for _, field := range f.Facets {
			f.Query.Add("attributesForFaceting", field.Attr())
		}
	}
	return f.Query.Encode()
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

func valsToMap(pm map[string]any, q url.Values) {
	for attr, vals := range q {
		switch attr {
		case "attributesForFaceting":
			if len(vals) == 1 {
				pm[attr] = strings.Split(vals[0], ",")
			} else {
				pm[attr] = vals
			}
		default:
			pm[attr] = vals
		}
	}
}

func GetDataFromQuery(q url.Values) []map[string]any {
	if q.Has("data") {
		d, err := FileSrc(q["data"])
		if err != nil {
			return []map[string]any{}
		}
		return d
	}
	return []map[string]any{}
}

// FileSrc takes json data files.
func FileSrc(files []string) ([]map[string]any, error) {
	var data []map[string]any
	for _, file := range files {
		p, err := dataFromFile(file)
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

// DecodeData decodes data from a io.Reader.
func DecodeData(r io.Reader) ([]map[string]any, error) {
	var data []map[string]any
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
