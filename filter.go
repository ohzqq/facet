package facet

import (
	"fmt"
	"net/url"

	"github.com/spf13/cast"
)

func Filter(data []map[string]any, cfg map[string]any, filters url.Values) ([]map[string]any, map[string]any) {
	idx := &Index{
		Data: data,
	}
	if f, ok := cfg["facets"]; ok {
		idx.FacetCfg = parseFacetMap(f)
	}
	if key, ok := cfg["key"]; ok {
		idx.Key = cast.ToString(key)
	}
	idx.Facets()
	fmt.Printf("%#v\n", idx)
	return data, cfg
}

func parseFilters(f any) (url.Values, error) {
	filters := make(map[string][]string)
	var err error
	switch val := f.(type) {
	case url.Values:
		return val, nil
	case []byte:
		filters, err = cast.ToStringMapStringSliceE(string(val))
		if err != nil {
			return nil, err
		}
	case string:
		q, err := url.ParseQuery(val)
		if err != nil {
			return nil, err
		}
		return q, nil
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	return url.Values(filters), nil
}
