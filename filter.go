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
