package facet

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Facet is a structure for facet data.
type Facet struct {
	Attribute string       `json:"attribute"`
	Items     []*FacetItem `json:"items,omitempty"`
	Operator  string       `json:"operator,omitempty"`
	Sep       string       `json:"-"`
}

// NewFacet initializes a facet with attribute.
func NewFacet(name string) *Facet {
	return &Facet{
		Attribute: name,
		Operator:  "or",
		Sep:       ".",
	}
}

// GetItem returns an *FacetItem.
func (f *Facet) GetItem(term string) *FacetItem {
	for _, item := range f.Items {
		if term == item.Value {
			return item
		}
	}
	return f.AddItem(term)
}

// GetConfig returns a map of a Facet's config.
func (f *Facet) GetConfig() map[string]any {
	return map[string]any{
		"attribute": f.Attribute,
		"operator":  f.Operator,
	}
}

// ListItems returns a string slice of all item values.
func (f *Facet) ListItems() []string {
	var items []string
	for _, item := range f.Items {
		items = append(items, item.Value)
	}
	return items
}

// AddItem adds an item with optional ids. If the item already exists ids are
// appended.
func (f *Facet) AddItem(term string, ids ...string) *FacetItem {
	for _, i := range f.Items {
		if term == i.Value {
			i.BelongsTo(ids...)
			return i
		}
	}
	item := NewFacetItem(term, ids)
	f.Items = append(f.Items, item)
	return item
}

// CollectItems takes the input data and aggregates them based on the
// Facet.Attribute.
func (f *Facet) CollectItems(data []map[string]any) *Facet {
	for i, item := range data {
		if terms, ok := item[f.Attribute]; ok {
			var items []string
			switch t := terms.(type) {
			case string:
				items = append(items, t)
			case []string:
				items = t
			case []any:
				items = cast.ToStringSlice(t)
			}
			for _, term := range items {
				f.AddItem(term, cast.ToString(i))
			}
		}
	}
	return f
}

// FuzzyFindItem fuzzy finds an item's value and returns possible matches.
func (f *Facet) FuzzyFindItem(term string) []*FacetItem {
	matches := f.FuzzyMatches(term)
	items := make([]*FacetItem, len(matches))
	for i, match := range matches {
		item := f.Items[match.Index]
		item.Match = match
		items[i] = item
	}
	return items
}

// FuzzyMatches returns the fuzzy.Matches of the search.
func (f *Facet) FuzzyMatches(term string) fuzzy.Matches {
	return fuzzy.FindFrom(term, f)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (f *Facet) String(i int) string {
	return f.Items[i].Value
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Facet) Len() int {
	return len(f.Items)
}

// Filter applies the listed filters to the facet.
func (f *Facet) Filter(filters ...string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, filter := range filters {
		term := f.FuzzyFindItem(filter)
		if len(term) < 1 {
		}
		bits = append(bits, term[0].Bitmap())
	}

	switch f.Operator {
	case "and":
		return roaring.ParAnd(viper.GetInt("workers"), bits...)
	default:
		return roaring.ParOr(viper.GetInt("workers"), bits...)
	}
}

// FacetItem is a data structure for a Facet's item.
type FacetItem struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Count       int    `json:"count"`
	belongsTo   []uint32
	fuzzy.Match `json:"-"`
}

// NewFacetItem initializes an item with a value and string slice of related data
// items.
func NewFacetItem(name string, vals []string) *FacetItem {
	term := &FacetItem{
		Value: name,
		Label: name,
	}
	term.BelongsTo(vals...)
	return term
}

// BelongsTo adds slice of index values for data items.
func (t *FacetItem) BelongsTo(vals ...string) *FacetItem {
	for _, val := range vals {
		t.belongsTo = append(t.belongsTo, cast.ToUint32(val))
	}
	t.Count = len(t.belongsTo)
	return t
}

// Bitmap returns a *roaring.Bitmap of slice indices for a FacetItem.
func (t *FacetItem) Bitmap() *roaring.Bitmap {
	return roaring.BitmapOf(t.belongsTo...)
}
