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
	idx = New("audiobooks", []string{"tags"}, books...)
	fmt.Printf("%v\n", idx.Name)
}

func TestProcessFacets(t *testing.T) {
	//idx = idx.ProcessFacets()
	idx.processData()
	fmt.Printf("%v\n", idx.getIDs())
	for _, f := range idx.facets {
		og := len(lo.Keys(f.terms))
		uniq := len(lo.Uniq(lo.Keys(f.terms)))
		if og != uniq {
			t.Errorf("got %d terms, expected %d\n", uniq, og)
		}
	}
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
