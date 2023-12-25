package facet

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/samber/lo"
)

var idx *Index

func init() {
	idx = New("audiobooks", []string{"tags", "authors", "narrators"}, books)
}

func TestNewIndex(t *testing.T) {
	books := loadData(t)
	idx = New("audiobooks", []string{"tags", "authors", "narrators"}, books)
	fmt.Printf("%v\n", idx.Name)
}

func TestIdxCfg(t *testing.T) {
	cfg := &Index{}
	err := json.Unmarshal([]byte(testCfg), cfg)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", cfg)
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

func TestConjQuery(t *testing.T) {
	abo := idx.GetFacetTermItems("tags", "abo")
	dnr := idx.GetFacetTermItems("tags", "dnr")
	fmt.Printf("abo %v\n", len(abo))
	fmt.Printf("dnr %v\n", len(dnr))
	or := lo.Union(abo, dnr)

	books := idx.GetByID(or)
	if len(or) != len(books) {
		t.Errorf("got %d books expected %d\n", len(books), len(or))
	}
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
