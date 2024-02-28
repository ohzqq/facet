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
	Keywords  []*Keyword
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
	tokens := make(map[string]map[string]any)
	for _, token := range f.Keywords {
		tokens[token.Label] = map[string]any{
			"count": token.Count(),
			"items": token.Items(),
		}
	}
	d, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (f *Field) FindByLabel(label string) *Keyword {
	for _, token := range f.Keywords {
		if token.Label == label {
			return token
		}
	}
	return NewKeyword(label)
}

func (f *Field) FindByValue(val string) *Keyword {
	for _, token := range f.Keywords {
		if token.Value == val {
			return token
		}
	}
	return NewKeyword(val)
}

func (f *Field) FindByIndex(ti ...int) []*Keyword {
	var tokens []*Keyword
	for _, tok := range ti {
		if tok < f.Count() {
			tokens = append(tokens, f.Keywords[tok])
		}
	}
	return tokens
}

func (f *Field) Add(val any, ids []int) {
	for _, token := range f.Tokenize(val) {
		if f.kwIdx == nil {
			f.kwIdx = make(map[string]int)
		}
		if idx, ok := f.kwIdx[token.Value]; ok {
			f.Keywords[idx].Add(ids...)
		} else {
			idx = len(f.Keywords)
			f.kwIdx[token.Value] = idx
			f.Keywords = append(f.Keywords, token)
		}
	}
}

func (f *Field) Tokenize(val any) []*Keyword {
	return f.analyzer.Tokenize(val)
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

func (f *Field) Search(term string) []*Keyword {
	matches := fuzzy.FindFrom(term, f)
	tokens := make([]*Keyword, len(matches))
	for i, match := range matches {
		tokens[i] = f.Keywords[match.Index]
	}
	return tokens
}

func (f *Field) Filter(val string) *roaring.Bitmap {
	tokens := f.Find(val)
	bits := make([]*roaring.Bitmap, len(tokens))
	for i, token := range tokens {
		bits[i] = token.Bitmap()
	}
	return roaring.ParAnd(viper.GetInt("workers"), bits...)
}

func (f *Field) Fuzzy(term string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(term, f)
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		b := f.Keywords[match.Index].Bitmap()
		bits[i] = b
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

func (f *Field) GetValues() []string {
	vals := make([]string, len(f.Keywords))
	for i, token := range f.Keywords {
		vals[i] = token.Value
	}
	return vals
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Field) Len() int {
	return f.Count()
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (f *Field) String(i int) string {
	return f.Keywords[i].Label
}

func (f *Field) Find(val any) []*Keyword {
	var tokens []*Keyword
	for _, tok := range f.Tokenize(val) {
		if token, ok := f.kwIdx[tok.Value]; ok {
			tokens = append(tokens, f.Keywords[token])
		}
	}
	return tokens
}

func (f *Field) Count() int {
	return len(f.Keywords)
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
