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
		idx.Facets = parseFacetMap(f)
	}
	fmt.Printf("%#v\n", idx)
	return data, cfg
}

func FilterItems(data []map[string]any, ids []any) []map[string]any {
	items := make([]map[string]any, len(ids))
	for item, _ := range data {
		for i, id := range ids {
			if cast.ToInt(id) == item {
				items[i] = data[item]
			}
		}
	}
	return items
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
