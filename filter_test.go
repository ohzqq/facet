package facet

import (
	"net/url"
	"testing"
)

var filterStrs = []filterStr{
	filterStr{
		want:  2241,
		query: `[j"tags:dnr"]`,
	},
	filterStr{
		want:  384,
		query: `["tags:dnr", "tags:abo"]`,
	},
	filterStr{
		want:  32,
		query: `["tags:-dnr", "tags:abo"]`,
	},
	filterStr{
		want:  32,
		query: `["tags:abo", "tags:-dnr"]`,
	},
	filterStr{
		want:  2273,
		query: `[["tags:dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5395,
		query: `[["tags:-dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  5395,
		query: `[[ "tags:abo", "tags:-dnr"]]`,
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

func TestFilterVals(t *testing.T) {
	data, err := loadData()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range testSearchFilterStrings() {
		filter := f.vals.Get("facetFilters")
		filters, err := ParseFilters(filter)
		if err != nil {
			t.Error(err)
		}
		facets := New(data, defFieldsSlice, "id", filters...)
		if num := facets.Len(); num != f.want {
			t.Errorf("got %d results, wanted %d\nfilters: %v\n", num, f.want, f.vals.Get("facetFilters"))
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
