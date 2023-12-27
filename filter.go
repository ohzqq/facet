package facet

import (
	"net/url"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func Filter(idx *Index) *Index {
	var bits []*roaring.Bitmap
	for name, filters := range idx.Filters {
		if facet, ok := idx.Facets[name]; ok {
			bits = append(bits, facet.Filter(filters...))
		}
	}

	filtered := roaring.ParOr(4, bits...)

	ids := filtered.ToArray()
	data := FilteredItems(idx.Data, lo.ToAnySlice(ids))
	res := &Index{
		Data:   data,
		Facets: idx.Facets,
	}
	return res.CollectTerms()
}

func FilteredItems(data []map[string]any, ids []any) []map[string]any {
	items := make([]map[string]any, len(ids))
	for item, _ := range data {
		for i, id := range ids {
			if cast.ToInt(id) == item {
				items[i] = data[item]
			}
		}
	}
	return items
}

func FilterString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func FilterBytes(val []byte) (url.Values, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return filters, nil
}

func FilterAny(f any) (url.Values, error) {
	return parseFilters(f)
}

func parseFilters(f any) (url.Values, error) {
	filters := make(map[string][]string)
	var err error
	switch val := f.(type) {
	case url.Values:
		return val, nil
	case []byte:
		return FilterBytes(val)
	case string:
		return FilterString(val)
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	return url.Values(filters), nil
}
