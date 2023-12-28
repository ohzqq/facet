package facet

import (
	"net/url"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Filter takes an *Index, filters the data and calculates the facets. It
// returns a new *Index.
func Filter(idx *Index) *Index {
	var bits []*roaring.Bitmap
	for name, filters := range idx.Filters {
		for _, facet := range idx.Facets {
			if facet.Attribute == name {
				bits = append(bits, facet.Filter(filters...))
			}
		}
	}

	filtered := roaring.ParOr(viper.GetInt("workers"), bits...)
	ids := filtered.ToArray()

	res, err := New(
		idx.GetConfig(),
		FilteredItems(idx.Data, lo.ToAnySlice(ids)),
	)
	res.Filters = idx.Filters
	if err != nil {
		return res
	}
	return res
}

// FilteredItems returns the subset of data.
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

// FilterString parses an encoded filter string.
func FilterString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// FilterBytes parses a byte slice to url.Values.
func FilterBytes(val []byte) (url.Values, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return filters, nil
}

// ParseFilters takes an interface{} and returns a url.Values.
func ParseFilters(f any) (url.Values, error) {
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
