package facet

import (
	"encoding/json"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/cast"
)

const testDataFile = `testdata/data-dir/audiobooks.json`
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

func TestNewFacets(t *testing.T) {
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}
	facets := New(data, defFieldsSlice, "id")
	t.Fatalf("%#v\n", facets)
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
