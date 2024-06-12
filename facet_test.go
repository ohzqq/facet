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

var (
	testFilterStr = []string{
		`["tags:dnr", "tags:abo"]`,
		`["authors:amy lane", ["tags:-cops"]]`,
	}
)

var defFieldsStr = `tags,authors,narrators,series`
var defFieldsSingle = []string{"tags,authors,narrators,series"}
var defFieldsSlice = []string{"tags", "authors", "narrators", "series"}

var testQueryVals = url.Values{
	"attributesForFaceting": defFieldsSingle,
	"data":                  []string{"testdata/ndbooks.json"},
}

var facetCount = map[string]int{
	"tags":      217,
	"authors":   1602,
	"series":    1722,
	"narrators": 1412,
}

func TestNewFacets(t *testing.T) {
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}
	facets := New(data, defFieldsSlice, "id")
	for name, facet := range facets.Fields {
		got := facet.Len()
		want := facetCount[name]
		if got != want {
			t.Errorf("facet %s got %#v, wanted %v\n", name, got, want)
		}
	}
}

func TestParseFilterFacets(t *testing.T) {
	for _, test := range testFilterStr {
		_, err := ParseFilters(test)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFilterFacets(t *testing.T) {
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}
	for i, test := range testFilterStr {
		filters, err := ParseFilters(test)
		if err != nil {
			t.Error(err)
		}
		facets := New(data, defFieldsSlice, "id", filters...)
		wants := map[string]int{
			"tags":      23,
			"authors":   56,
			"series":    95,
			"narrators": 112,
		}
		for name, facet := range facets.Fields {
			if i == 0 {
				got := facet.Len()
				want := wants[name]
				if got != want {
					t.Errorf("facet %s got %#v, wanted %v\n", name, got, want)
				}
			}
		}
	}
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
