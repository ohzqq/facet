package facet

import (
	"net/url"
	"testing"
)

const (
	// search params
	Hits                 = `hits`
	AttributesToRetrieve = `attributesToRetrieve`
	Page                 = "page"
	HitsPerPage          = "hitsPerPage"
	SortFacetsBy         = `sortFacetValuesBy`
	ParamFacets          = "facets"
	ParamFilters         = "filters"
	FacetFilters         = `facetFilters`
	ParamFullText        = `fullText`
	NbHits               = `nbHits`
	NbPages              = `nbPages`
	SortBy               = `sortBy`
	Order                = `order`

	// Settings
	SrchAttr     = `searchableAttributes`
	FacetAttr    = `attributesForFaceting`
	SortAttr     = `sortableAttributes`
	DataDir      = `dataDir`
	DataFile     = `dataFile`
	DefaultField = `title`

	TextAnalyzer    = "text"
	KeywordAnalyzer = "keyword"
)

var filterStrs = []filterStr{
	filterStr{
		want:  2237,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr"]`,
	},
	filterStr{
		want:  384,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr", "tags:abo"]`,
	},
	filterStr{
		want:  2270,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  417,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["-tags:dnr", "tags:abo"]`,
	},
	filterStr{
		want:  417,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["-tags:dnr", "tags:abo"]]`,
	},
	filterStr{
		want:  2237,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr","-tags:abo"]`,
	},
	filterStr{
		want:  2237,
		query: `data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:dnr", "-tags:abo"]]`,
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
		if num := facets.Count(); num != f.want {
			t.Errorf("got %d results, wanted %d\n", num, f.want)
		}

		for _, facet := range facets.Facets {
			if num, ok := facetCount[facet.Attribute]; ok {
				if num != facet.Count() {
					t.Errorf("%v got %d, expected %d \n", facet.Attribute, facet.Count(), num)
				}
			} else {
				t.Errorf("attr %s not found\n", facet.Attribute)
			}
		}
	}

}

func TestFilterVals(t *testing.T) {
	for _, f := range testSearchFilterStrings() {
		facets, err := New(f.vals)
		if err != nil {
			t.Fatal(err)
		}
		if num := facets.Count(); num != f.want {
			t.Errorf("got %d results, wanted %d\n", num, f.want)
		}

		facets, err = New(f.vals.Encode())
		if err != nil {
			t.Fatal(err)
		}
		facets.Calculate()
		filtered, err := Filter(facets.bits, facets.Facets, facets.Filters())
		if err != nil {
			t.Fatal(err)
		}
		if num := filtered.GetCardinality(); num != uint64(f.want) {
			t.Errorf("got %d results, wanted %d\n", num, f.want)
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
			FacetFilters: []string{
				`["authors:amy lane"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 801,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 784,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 384,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["tags:dnr", "tags:abo"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 1853,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["tags:dnr", "tags:-abo"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 1853,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["tags:-abo", "tags:dnr"]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 2270,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`[["tags:dnr", "tags:abo"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 2237,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`[["tags:dnr", "tags:-abo"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 2237,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`[["tags:-abo", "tags:dnr"]]`,
			},
		},
	})

	queries = append(queries, filterVal{
		want: 0,
		vals: url.Values{
			"data":                  []string{"testdata/ndbooks.json"},
			"attributesForFaceting": []string{"tags", "authors"},
			FacetFilters: []string{
				`["tags:abo", "tags:dnr", "tags:horror"]`,
			},
		},
	})

	return queries
}
