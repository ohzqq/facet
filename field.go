package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/viper"
)

const (
	SortByCount = `count`
	SortByAlpha = `alpha`
)

type Field struct {
	Attribute string `json:"attribute"`
	Sep       string `json:"-"`
	SortBy    string
	Order     string
	keywords  []*Token
	kwIdx     map[string]int
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:    "/",
		SortBy: "count",
		Order:  "desc",
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
	field := make(map[string]any)
	field["facetValues"] = f.Keywords()
	if f.Len() < 1 {
		field["facetValues"] = []any{}
	}
	field["attribute"] = joinAttr(f)
	field["count"] = f.Len()
	return json.Marshal(field)
}

func (f *Field) Keywords() []*Token {
	return f.SortTokens()
}

func (f *Field) GetValues() []string {
	vals := make([]string, f.Len())
	for i, token := range f.keywords {
		vals[i] = token.Value
	}
	return vals
}

func (f *Field) FindByLabel(label string) *Token {
	for _, token := range f.keywords {
		if token.Label == label {
			return token
		}
	}
	return NewToken(label)
}

func (f *Field) FindByValue(val string) *Token {
	for _, token := range f.keywords {
		if token.Value == val {
			return token
		}
	}
	return NewToken(val)
}

func (f *Field) FindByIndex(ti ...int) []*Token {
	var tokens []*Token
	for _, tok := range ti {
		if tok < f.Len() {
			tokens = append(tokens, f.keywords[tok])
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
			f.keywords[idx].Add(ids...)
		} else {
			idx = len(f.keywords)
			f.kwIdx[token.Value] = idx
			token.Add(ids...)
			f.keywords = append(f.keywords, token)
		}
	}
}

func (f *Field) Tokenize(val any) []*Token {
	return KeywordTokenizer(val)
}

func (f *Field) Search(term string) []*Token {
	matches := fuzzy.FindFrom(term, f)
	tokens := make([]*Token, len(matches))
	for i, match := range matches {
		tokens[i] = f.keywords[match.Index]
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

func (f *Field) Find(val any) []*Token {
	var tokens []*Token
	for _, tok := range f.Tokenize(val) {
		if token, ok := f.kwIdx[tok.Value]; ok {
			tokens = append(tokens, f.keywords[token])
		}
	}
	return tokens
}

func (f *Field) Fuzzy(term string) *roaring.Bitmap {
	matches := fuzzy.FindFrom(term, f)
	bits := make([]*roaring.Bitmap, len(matches))
	for i, match := range matches {
		b := f.keywords[match.Index].Bitmap()
		bits[i] = b
	}
	return roaring.ParOr(viper.GetInt("workers"), bits...)
}

// Len returns the number of items, to satisfy the fuzzy.Source interface.
func (f *Field) Len() int {
	return len(f.keywords)
}

// String returns an Item.Value, to satisfy the fuzzy.Source interface.
func (f *Field) String(i int) string {
	return f.keywords[i].Label
}

func joinAttr(field *Field) string {
	attr := field.Attribute
	if field.SortBy != "" {
		attr += ":"
		attr += field.SortBy
	}
	if field.Order != "" {
		attr += ":"
		attr += field.Order
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
