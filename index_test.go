package facet

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var idx *Index

func TestNewIndex(t *testing.T) {
	books := loadData(t)
	idx = New("audiobooks", []string{"tags"}, books...)
	fmt.Printf("%v\n", idx.Name)
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
