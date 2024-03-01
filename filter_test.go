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

var filterStrs = []string{
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr"]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr", "tags:abo"]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:dnr", "tags:abo"]]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["-tags:dnr", "tags:abo"]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["-tags:dnr", "tags:abo"]]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=["tags:dnr","-tags:abo"]`,
	`data=testdata/ndbooks.json&attributesForFaceting=tags&facetFilters=[["tags:dnr", "-tags:abo"]]`,
}

type filterStr struct {
	vals url.Values
	want int
}

func TestFilters(t *testing.T) {
	for _, query := range filterStrs {
		facets, err := New(query)
		if err != nil {
			t.Fatal(err)
		}
		facets.Calculate()
		filtered, err := Filter(facets.bits, facets.Facets, facets.Filters())
		if err != nil {
			t.Fatal(err)
		}
		println(filtered.GetCardinality())
	}

}

func testSearchFilterStrings() []filterStr {
	//queries := make(map[int]url.Values)
	var queries []filterStr

	queries = append(queries, filterStr{
		want: 58,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 806,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 789,
		vals: url.Values{
			FacetFilters: []string{
				`["authors:amy lane", ["tags:romance", "tags:-dnr"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 384,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:dnr", "tags:abo"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 1856,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:dnr", "tags:-abo"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 1856,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:-abo", "tags:dnr"]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2273,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2240,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:dnr", "tags:-abo"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 2240,
		vals: url.Values{
			FacetFilters: []string{
				`[["tags:-abo", "tags:dnr"]]`,
			},
		},
	})

	queries = append(queries, filterStr{
		want: 0,
		vals: url.Values{
			FacetFilters: []string{
				`["tags:abo", "tags:dnr", "tags:horror"]`,
			},
		},
	})

	return queries
}
