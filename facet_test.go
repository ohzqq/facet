package facet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/cast"
)

const testDataFile = `testdata/data-dir/ndbooks.json`
const testDataDir = `testdata/data-dir`
const numBooks = 7252
const testQueryString = `attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=narrators&attributesForFaceting=series&data=testdata/ndbooks.json&uid=url`

var queryStrTests = []string{
	`attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=narrators&attributesForFaceting=series&data=testdata/ndbooks.json&uid=id&facetFilters=["tags:dnr", "tags:abo"]`,
	`attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=narrators&attributesForFaceting=series&data=testdata/ndbooks.json&facetFilters=["tags:dnr", "tags:abo"]`,
}

var defFieldsStr = `tags,authors,narrators,series`
var defFieldsSingle = []string{"tags,authors,narrators,series"}
var defFieldsSlice = []string{"tags", "authors", "narrators", "series"}

var testQueryVals = url.Values{
	"attributesForFaceting": defFieldsSingle,
	"data":                  []string{"testdata/ndbooks.json"},
}

var facetCount = map[string]int{
	"tags":      218,
	"authors":   1612,
	"series":    1740,
	"narrators": 1428,
}

func TestNewFacetsFromQueryString(t *testing.T) {
	facets, err := New(testQueryString)
	if err != nil {
		t.Fatal(err)
	}

	err = testFacetCfg(facets)
	if err != nil {
		t.Error(err)
	}

	if len(facets.data) != numBooks {
		t.Errorf("got %d items, expected %d\n", len(facets.data), 7174)
	}
	//if len(facets.Hits) > 0 {
	//  fmt.Printf("%+v\n", facets.Hits[0]["title"])
	//}
}

func TestNewFacetsFromQuery(t *testing.T) {
	facets, err := New(testQueryVals)
	if err != nil {
		t.Fatal(err)
	}

	err = testFacetCfg(facets)
	if err != nil {
		t.Error(err)
	}

	for _, facet := range facets.Facets {
		if num, ok := facetCount[facet.Attribute]; ok {
			if num != facet.Len() {
				t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Len(), num)
			}
		} else {
			t.Errorf("attr %s not found\n", facet.Attribute)
		}
	}
}

func testFacetCfg(facets *Facets) error {
	if attrs := facets.Attrs(); len(attrs) != 4 {
		return fmt.Errorf("got %d attributes, expected %d\n", len(attrs), 4)
	}

	facets.Calculate()
	if len(facets.Facets) != 4 {
		return fmt.Errorf("got %d attributes, expected %d\n", len(facets.Facets), 4)
	}

	return nil
}

func dataToMap() (map[string]map[string]any, error) {
	data, err := loadData()
	if err != nil {
		return nil, err
	}

	d := make(map[string]map[string]any)
	for _, i := range data {
		id := cast.ToString(i["id"])
		d[id] = i
	}
	return d, nil
}

func loadData() ([]map[string]any, error) {
	d, err := os.ReadFile(testDataFile)
	if err != nil {
		return nil, err
	}

	var books []map[string]any
	err = json.Unmarshal(d, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}
