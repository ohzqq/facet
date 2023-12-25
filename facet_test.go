package facet

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

var books []map[string]any

const numBooks = 7174

func init() {
	d, err := os.ReadFile("testdata/audiobooks.json")
	if err != nil {
		log.Fatal(err)
	}

	var res []map[string]any
	err = json.Unmarshal(d, &res)
	if err != nil {
		log.Fatal(err)
	}

	books = res
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
