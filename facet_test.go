package facet

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cast"
)

const testDataFile = `testdata/data-dir/audiobooks.json`
const testDataDir = `testdata/data-dir`
const numBooks = 7253

func TestFacets(t *testing.T) {
	data, err := dataToMap()
	if err != nil {
		t.Fatal(err)
	}

	facets := New(data, []string{"tags", "authors", "narrators", "series"})
	for _, facet := range facets.fields {
		fmt.Printf("%+v\n", facet.Count())
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
