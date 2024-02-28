package facet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/cast"
)

const testDataFile = `testdata/data-dir/audiobooks.json`
const testDataDir = `testdata/data-dir`
const numBooks = 7253
const testQueryString = `attributesForFaceting=tags&attributesForFaceting=authors&attributesForFaceting=narrators&attributesForFaceting=series&data=testdata/data-dir/audiobooks.json&data=testdata/audiobooks.json&&uid=id`

var defFieldsStr = `tags,authors,narrators,series`
var defFieldsSingle = []string{"tags,authors,narrators,series"}
var defFieldsSlice = []string{"tags", "authors", "narrators", "series"}

var testQueryVals = url.Values{
	"attributesForFaceting": defFieldsSingle,
	"data":                  []string{"testdata/audiobooks.json"},
}

func TestFacets(t *testing.T) {
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}

	facets := NewFacets(data, []string{"tags", "authors", "narrators", "series"})
	facets.Calculate()
	for _, facet := range facets.Facets {
		fmt.Printf("%+v\n", facet.Count())
	}
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

	facets.Calculate()
	for _, facet := range facets.Facets {
		fmt.Printf("%+v\n", facet.Count())
	}

	//if len(facets.Data)
	println(len(facets.Data))
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
}

func TestNewFacetsFromJSON(t *testing.T) {
	raw := make(map[string]any)
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}
	raw["data"] = data
	raw["attributesForFaceting"] = defFieldsSlice

	facets, err := New(raw)
	if err != nil {
		t.Fatal(err)
	}

	err = testFacetCfg(facets)
	if err != nil {
		t.Error(err)
	}

	d, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}

	facets, err = New(d)
	if err != nil {
		t.Fatal(err)
	}

	err = testFacetCfg(facets)
	if err != nil {
		t.Error(err)
	}
}

func testFacetCfg(facets *Facets) error {
	if len(facets.Attrs) != 4 {
		return fmt.Errorf("got %d attributes, expected %d\n", len(facets.Attrs), 4)
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
