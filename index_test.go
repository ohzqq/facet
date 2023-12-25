package facet

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/samber/lo"
)

var idx = &Index{}

var books []map[string]any

const numBooks = 7174

func init() {
	cfg, err := os.ReadFile("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cfg, idx)
	if err != nil {
		log.Fatal(err)
	}

	d, err := os.ReadFile("testdata/audiobooks.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(d, &idx.Data)
	if err != nil {
		log.Fatal(err)
	}

	books = idx.Data
}

func TestNewIndex(t *testing.T) {
	New("audiobooks", []string{"tags", "authors", "narrators"}, books)
}

func TestIdxCfg(t *testing.T) {
	//cfg := &Index{}
	err := json.Unmarshal([]byte(testCfg), idx)
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("%v\n", cfg)
}

func TestProcessFacets(t *testing.T) {
	//idx = idx.ProcessFacets()

	//idx.processData()
	for _, f := range idx.facets {
		og := len(lo.Keys(f))
		uniq := len(lo.Uniq(lo.Keys(f)))
		if og != uniq {
			t.Errorf("got %d terms, expected %d\n", uniq, og)
		}
		//fmt.Printf("%v\n", f)
	}

	//terms := idx.FacetMap()
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
