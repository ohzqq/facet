package facet

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7174

func init() {
	var err error
	idx, err = New("testdata/config.json", "testdata/audiobooks.json")
	if err != nil {
		log.Fatal(err)
	}
	books = idx.Data
}

func TestIdxCfg(t *testing.T) {
	//cfg := &Index{}
	err := json.Unmarshal([]byte(testCfg), idx)
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("%v\n", cfg)
}

func loadData(t *testing.T) []map[string]any {
	d, err := os.ReadFile("testdata/audiobooks.json")
	if err != nil {
		t.Error(err)
	}

	var books []map[string]any
	err = json.Unmarshal(d, &books)
	if err != nil {
		t.Error(err)
	}

	books = books

	return books
}

func TestData(t *testing.T) {
	books := loadData(t)
	println(len(books))
}

const testCfg = `
{
	"name": "audiobooks",
	"key": "id",
	"facets": {
		"tags": {
			"operator": "and"
		},
		"authors": {}
	}
}
`
