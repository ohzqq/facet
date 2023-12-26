package facet

import (
	"net/url"
	"testing"
)

func TestFilters(t *testing.T) {
	println("test filters")
	f1 := `tags=abo&tags=dnr&authors=Alice+Winters&authors=Amy+Lane`
	q, err := url.ParseQuery(f1)
	if err != nil {
		t.Error(err)
	}
	testFilters(q)
	//fmt.Printf("%+v\n", books)
}
