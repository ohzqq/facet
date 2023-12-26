package facet

import (
	"fmt"
	"log"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

func Filter(data []map[string]any, f map[string]map[string]any, filters url.Values) ([]map[string]any, map[string]map[string]any) {
	facets := make(map[string]*Facet)
	for name, agg := range f {
		facet := NewFacet(name)
		err := mapstructure.Decode(agg, facet)
		if err != nil {
			log.Fatal(err)
		}
		facets[f] = facet
	}
	fmt.Printf("%#v\n", facets)
}
