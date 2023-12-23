package facet

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/samber/lo"
)

var idx *Index

func TestNewIndex(t *testing.T) {
	books := loadData(t)
	idx = New("audiobooks", []string{"tags", "authors", "narrators"}, books)
	fmt.Printf("%v\n", idx.Name)
}

func TestProcessFacets(t *testing.T) {
	//idx = idx.ProcessFacets()

	//idx.processData()
	for _, f := range idx.Facets {
		og := len(lo.Keys(f))
		uniq := len(lo.Uniq(lo.Keys(f)))
		if og != uniq {
			t.Errorf("got %d terms, expected %d\n", uniq, og)
		}
		//fmt.Printf("%v\n", f)
	}

	//terms := idx.FacetMap()
	fmt.Printf("%v\n", idx.GetFacetValues("tags"))
}

func TestData(t *testing.T) {
	books := loadData(t)
	println(len(books))
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

	return books
}
