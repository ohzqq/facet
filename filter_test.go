package facet

import (
	"encoding/json"
	"net/url"
	"os"
	"testing"
)

var filterStrs = []filterStr{
	filterStr{
		want:  2241,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr"]`,
	},
	filterStr{
		want:  384,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr", "tags:abo"]`,
	},
	filterStr{
		want:  32,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:-dnr", "tags:abo"]`,
	},
	filterStr{
		want:  32,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:abo", "tags:-dnr"]`,
	},
	filterStr{
		want:  2273,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5395,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:-dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5395,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[[ "tags:abo", "tags:-dnr"]]`,
	},
}

type filterStr struct {
	query string
	want  int
}

type filterVal struct {
	vals url.Values
	want int
}

func TestFilterStrings(t *testing.T) {
	for _, f := range filterStrs {
		facets, err := New(f.query)
		if err != nil {
			t.Fatal(err)
		}
		if num := facets.Len(); num != f.want {
			t.Errorf("query %s:\ngot %d results, wanted %d\n", f.query, num, f.want)
		}
	}

}

func TestFilterVals(t *testing.T) {
	for _, f := range testSearchFilterStrings() {
		facets, err := New(f.vals)
		if err != nil {
			t.Fatal(err)
		}
		if num := facets.Len(); num != f.want {
			t.Errorf("got %d results, wanted %d\nfilters: %v\n", num, f.want, f.vals.Get("facetFilters"))
		}
		//if len(facets.Hits) > 0 {
		//  fmt.Printf("%v: %+v\n", f.vals.Encode(), facets.Hits[0]["title"])
		//}
		//println(facets.Len())

		facets, err = New(f.vals.Encode())
		if err != nil {
			t.Fatal(err)
		}
		if num := facets.Len(); num != f.want {
			t.Errorf("got %d results, wanted %d\nfilters: %v\n", num, f.want, f.vals.Get("facetFilters"))
		}

		enc, err := json.Marshal(facets)
		if err != nil {
			t.Error(err)
		}
		println(string(enc))

		err = facets.Encode(os.Stdout)
		if err != nil {
			t.Error(err)
		}
	}
}

func testSearchFilterStrings() []filterVal {
	//queries := make(map[int]url.Values)
	var queries []filterVal

	queries = append(queries, filterVal{
		want: 58,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 26,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 41,
		vals: url.Values{
			//"uid": []string{"id"},
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 0,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			"facetFilters": []string{
				`["tags:abo", "tags:dnr", "tags:horror"]`,
			},
		},
	})

	return queries
}
