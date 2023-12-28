package facet

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Facet struct {
	Attribute string  `json:"attribute"`
	Items     []*Term `json:"items,omitempty"`
	Operator  string  `json:"operator,omitempty"`
	Sep       string  `json:"-"`
}

func NewFacet(name string) *Facet {
	return &Facet{
		Attribute: name,
		Operator:  "or",
		Sep:       ".",
	}
}

func (f *Facet) GetItem(term string) *Term {
	for _, item := range f.Items {
		if term == item.Value {
			return item
		}
	}
	return f.AddItem(term)
}

func (f *Facet) ListItems() []string {
	var items []string
	for _, item := range f.Items {
		items = append(items, item.Value)
	}
	return items
}

func (f *Facet) AddItem(term string, ids ...string) *Term {
	for _, i := range f.Items {
		if term == i.Value {
			i.AddItems(ids...)
			return i
		}
	}
	item := NewTerm(term, ids)
	f.Items = append(f.Items, item)
	return item
}

func (f *Facet) CollectItems(data []map[string]any) *Facet {
	for i, item := range data {
		if terms, ok := item[f.Attribute]; ok {
			for _, term := range terms.([]any) {
				f.AddItem(cast.ToString(term), cast.ToString(i))
			}
		}
	}
	return f
}

func (f *Facet) HierarchicalItems(sep string) []*Term {
	fn := func(item *Term, i int) bool {
		b := strings.Contains(item.Value, sep)
		return b
	}
	items := lo.Filter(f.Items, fn)
	return items
}

func (f *Facet) FuzzyFindItem(term string) []*Term {
	matches := f.FuzzyMatches(term)
	items := make([]*Term, len(matches))
	for i, match := range matches {
		item := f.Items[match.Index]
		item.Match = match
		items[i] = item
	}
	return items
}

func (f *Facet) FuzzyMatches(term string) fuzzy.Matches {
	return fuzzy.FindFrom(term, f)
}

func (f *Facet) String(i int) string {
	return f.Items[i].Value
}

func (f *Facet) Len() int {
	return len(f.Items)
}

func (f *Facet) Filter(filters ...string) *roaring.Bitmap {
	var bits []*roaring.Bitmap
	for _, filter := range filters {
		term := f.GetItem(filter)
		bits = append(bits, term.Bitmap())
	}

	switch f.Operator {
	case "and":
		return roaring.ParAnd(4, bits...)
	default:
		return roaring.ParOr(4, bits...)
	}
}

type Term struct {
	Value       string  `json:"value"`
	Label       string  `json:"label"`
	Count       int     `json:"count"`
	Items       []*Term `json:"-"`
	belongsTo   []uint32
	fuzzy.Match `json:"-"`
}

func NewTerm(name string, vals []string) *Term {
	term := &Term{
		Value: name,
		Label: name,
	}
	term.AddItems(vals...)
	return term
}

func TermIsHierarchical(name, sep string) bool {
	return strings.Contains(name, sep)
}

func (t *Term) AddItems(vals ...string) *Term {
	for _, val := range vals {
		t.belongsTo = append(t.belongsTo, cast.ToUint32(val))
	}
	t.Count = len(t.belongsTo)
	return t
}

func (t *Term) HasChildren(sep ...string) bool {
	s := "."
	if len(sep) > 0 {
		s = sep[0]
	}
	return strings.Contains(t.Value, s) || len(t.Items) > 0
}

func (t *Term) Bitmap() *roaring.Bitmap {
	return roaring.BitmapOf(t.belongsTo...)
}
