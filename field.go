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
	terms     []string
	Keywords  map[string]*Keyword
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
	for _, label := range f.terms {
		token := f.FindByLabel(label)
		tokens[label] = token.Count()
	}
	d, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (t *Field) GetTokens() []*Keyword {
	var tokens []*Keyword
	for _, label := range t.terms {
		tok := t.FindByLabel(label)
		tokens = append(tokens, tok)
	}
	return tokens
}

func (t *Field) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		if t.Keywords == nil {
			t.Keywords = make(map[string]*Keyword)
		}
		if _, ok := t.Keywords[token.Value]; !ok {
			t.terms = append(t.terms, token.Label)
			t.Keywords[token.Value] = token
		}
		t.Keywords[token.Value].Add(ids...)
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

func (t *Field) FindByIndex(ti ...int) []*Keyword {
	var tokens []*Keyword
	toks := t.GetTokens()
	total := t.Count()
	for _, tok := range ti {
		if tok < total {
			tokens = append(tokens, toks[tok])
		}
	}
	return tokens
}

func (t *Field) Search(term string) []*Keyword {
	matches := fuzzy.FindFrom(term, t)
	tokens := make([]*Keyword, len(matches))
	all := t.GetTokens()
	for i, match := range matches {
		tokens[i] = all[match.Index]
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
	all := t.GetTokens()
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		b := all[match.Index].Bitmap()
		bits[i] = b
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

func (t *Field) GetValues() []string {
	sorted := t.GetTokens()
	tokens := make([]string, len(sorted))
	for i, t := range sorted {
		tokens[i] = t.Value
	}
	return tokens
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (t *Field) Len() int {
	return t.Count()
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (t *Field) String(i int) string {
	return t.terms[i]
}

func (t *Field) Find(val any) []*Keyword {
	var tokens []*Keyword
	for _, tok := range t.Tokenize(val) {
		if token, ok := t.Keywords[tok.Value]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (t *Field) FindByLabel(label string) *Keyword {
	for _, token := range t.Keywords {
		if token.Label == label {
			return token
		}
	}
	return NewItem(label)
}

func (t *Field) Count() int {
	return len(t.Keywords)
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
