package facet

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestFilters(t *testing.T) {
	d, err := os.ReadFile("testdata/config.json")
	if err != nil {
		t.Error(err)
	}
	tf := make(map[string]any)
	err = json.Unmarshal(d, &tf)
	if err != nil {
		t.Error(err)
	}
	facets := make(map[string]*Facet)
	if f, ok := tf["facets"]; ok {
		facets = parseFacetMap(f)
	} else {
		t.Errorf("facets not found in config")
	}

	fmt.Printf("%+v\n", facets)
	//fmt.Printf("%+v\n", books)
}
