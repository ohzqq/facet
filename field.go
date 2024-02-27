package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch/txt"
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
	*txt.Tokens
}

func NewField(attr string) *Field {
	f := &Field{
		Sep:    ".",
		Tokens: txt.NewTokens(),
	}
	f.SetAnalyzer(txt.Keyword())
	parseAttr(f, attr)
	return f
}

func NewFields(attrs []string) map[string]*Field {
	fields := make(map[string]*Field)
	for _, attr := range attrs {
		fields[attr] = NewField(attr)
	}

	return fields
}

func CalculateFacets(data []map[string]any, fields []string) map[string]*Field {
	facets := NewFields(fields)
	for id, d := range data {
		for attr, facet := range facets {
			if val, ok := d[attr]; ok {
				facet.Add(val, []int{id})
			}
		}
	}
	return facets
}

func (f *Field) MarshalJSON() ([]byte, error) {
	tokens := make(map[string]int)
	for _, label := range f.Tokens.Tokens {
		token := f.FindByLabel(label)
		tokens[label] = token.Count()
	}
	d, err := json.Marshal(tokens)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (t *Field) GetTokens() []*txt.Token {
	var tokens []*txt.Token
	for _, label := range t.Tokens.Tokens {
		tok := t.FindByLabel(label)
		tokens = append(tokens, tok)
	}
	return tokens
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

func (t *Field) FindByIndex(ti ...int) []*txt.Token {
	var tokens []*txt.Token
	toks := t.GetTokens()
	total := t.Count()
	for _, tok := range ti {
		if tok < total {
			tokens = append(tokens, toks[tok])
		}
	}
	return tokens
}

func (t *Field) Search(term string) []*txt.Token {
	matches := fuzzy.FindFrom(term, t)
	tokens := make([]*txt.Token, len(matches))
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
	return t.Tokens.Tokens[i]
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
