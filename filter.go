package facet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/mitchellh/mapstructure"
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

func parseCfg(c any) (string, map[string]*Facet, error) {
	var key string
	facets := make(map[string]*Facet)
	switch cfg := c.(type) {
	case []byte:
		return unmarshalCfg(cfg)
	case string:
		if Exist(cfg) {
			data, err := os.ReadFile(cfg)
			if err != nil {
				return "", nil, err
			}
			return unmarshalCfg(data)
		} else {
			return unmarshalCfg([]byte(cfg))
		}
	case map[string]any:
		return handleCfg(cfg)
	}

	return key, facets, nil
}

func unmarshalCfg(d []byte) (string, map[string]*Facet, error) {
	cfg := make(map[string]any)
	err := json.Unmarshal(d, &cfg)
	if err != nil {
		return "", nil, err
	}

	return handleCfg(cfg)
}

func handleCfg(cfg map[string]any) (string, map[string]*Facet, error) {
	pk := "id"

	facets := make(map[string]*Facet)
	if f, ok := cfg["facets"]; ok {
		facets = parseFacetMap(f)
	} else {
		return pk, facets, errors.New("facets not found in config")
	}

	if key, ok := cfg["key"]; ok {
		pk = cast.ToString(key)
	}

	return pk, facets, nil
}

func unmarshalFacets(d []byte) (map[string]*Facet, error) {
	facets := make(map[string]*Facet)
	err := json.Unmarshal(d, &facets)
	if err != nil {
		return nil, err
	}
	return facets, nil
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
