package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/viper"
)

const (
	Or               = "OR"
	And              = "AND"
	Not              = `NOT`
	AndNot           = `AND NOT`
	OrNot            = `OR NOT`
	FacetField       = `facet`
	SortByCount      = `count`
	SortByAlpha      = `alpha`
	StandardAnalyzer = `standard`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	keywords  []*Keyword
	kwIdx     map[string]int
	analyzer  Analyzer
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:      "/",
		analyzer: keyword{},
	}
	parseAttr(f, attr)
	return f
}

func NewFields(attrs []string) []*Field {
	fields := make([]*Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = NewField(attr)
	}
	return fields
}

func (f *Field) MarshalJSON() ([]byte, error) {
	tokens := make(map[string]int)
	for _, token := range f.keywords {
		//token := f.FindByLabel(label)
		tokens[token.Label] = token.Count()
	}
	d, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (t *Field) FindByLabel(label string) *Keyword {
	for _, token := range t.keywords {
		if token.Label == label {
			return token
		}
	}
	return NewItem(label)
}

func (f *Field) FindByIndex(ti ...int) []*Keyword {
	var tokens []*Keyword
	for _, tok := range ti {
		if tok < f.Count() {
			tokens = append(tokens, f.keywords[tok])
		}
	}
	return tokens
}

func (t *Field) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		if t.kwIdx == nil {
			t.kwIdx = make(map[string]int)
		}
		if idx, ok := t.kwIdx[token.Value]; ok {
			t.keywords[idx].Add(ids...)
		} else {
			idx = len(t.keywords)
			//t.labels = append(t.labels, token.Label)
			t.kwIdx[token.Value] = idx
			t.keywords = append(t.keywords, token)
		}
	}
}

func (t *Field) Tokenize(val any) []*Keyword {
	return t.analyzer.Tokenize(val)
}

func GetFieldItems(data []map[string]any, field *Field) []map[string]any {
	field.SortBy = SortByAlpha
	tokens := field.SortTokens()

	items := make([]map[string]any, len(tokens))
	for i, token := range tokens {
		items[i] = map[string]any{
			"attribute": field.Attribute,
			"value":     token.Value,
			"label":     token.Label,
			"count":     token.Count(),
			"hits":      ItemsByBitmap(data, token.Bitmap()),
		}
	}
	return items
}

func ItemsByBitmap(data []map[string]any, bits *roaring.Bitmap) []map[string]any {
	var res []map[string]any
	bits.Iterate(func(x uint32) bool {
		res = append(res, data[int(x)])
		return true
	})
	return res
}

func (t *Field) Search(term string) []*Keyword {
	matches := fuzzy.FindFrom(term, t)
	tokens := make([]*Keyword, len(matches))
	for i, match := range matches {
		tokens[i] = t.keywords[match.Index]
	}
	return tokens
}

func (t *Field) Filter(val string) *roaring.Bitmap {
	tokens := t.Find(val)
	bits := make([]*roaring.Bitmap, len(tokens))
	for i, token := range tokens {
		bits[i] = token.Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (t *Field) Fuzzy(term string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(term, t)
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		b := t.keywords[match.Index].Bitmap()
		bits[i] = b
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

func (t *Field) GetValues() []string {
	vals := make([]string, len(t.keywords))
	for i, token := range t.keywords {
		vals[i] = token.Value
	}
	return vals
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (t *Field) Len() int {
	return t.Count()
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (t *Field) String(i int) string {
	return t.keywords[i].Label
}

func (t *Field) Find(val any) []*Keyword {
	var tokens []*Keyword
	for _, tok := range t.Tokenize(val) {
		if token, ok := t.kwIdx[tok.Value]; ok {
			tokens = append(tokens, t.keywords[token])
		}
	}
	return tokens
}

func (t *Field) Count() int {
	return len(t.keywords)
}

func (f *Field) Attr() string {
	attr := f.Attribute
	if f.SortBy != "" {
		attr += ":"
		attr += f.SortBy
	}
	if f.Order != "" {
		attr += ":"
		attr += f.Order
	}
	return attr
}

func parseAttr(field *Field, attr string) {
	i := 0
	for attr != "" {
		var a string
		a, attr, _ = strings.Cut(attr, ":")
		if a == "" {
			continue
		}
		switch i {
		case 0:
			field.Attribute = a
		case 1:
			field.SortBy = a
		case 2:
			field.Order = a
		}
		i++
	}
}
