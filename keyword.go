package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Analyzer interface {
	Tokenize(any) []*Keyword
}

type Keyword struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Children *Field
	bits     *roaring.Bitmap
}

func NewKeyword(label string) *Keyword {
	return &Keyword{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Keyword) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Keyword) SetValue(txt string) *Keyword {
	f.Value = txt
	return f
}

func (f *Keyword) Count() int {
	return int(f.bits.GetCardinality())
}

func (f *Keyword) Contains(id int) bool {
	return f.bits.ContainsInt(id)
}

func (f *Keyword) Add(ids ...int) {
	for _, id := range ids {
		if !f.Contains(id) {
			f.bits.AddInt(id)
		}
	}
}

func (f *Keyword) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		f.Label: f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type keyword struct{}

func (kw keyword) Tokenize(str any) []*Keyword {
	return KeywordTokenizer(str)
}

func (kw keyword) Search(text string) []*Keyword {
	return []*Keyword{NewKeyword(normalizeText(text))}
}

func KeywordTokenizer(val any) []*Keyword {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Keyword, len(tokens))
	for i, token := range tokens {
		items[i] = NewKeyword(token)
		items[i].Value = normalizeText(token)
	}
	return items
}

func normalizeText(token string) string {
	fields := lowerCase(strings.Split(token, " "))
	for t, term := range fields {
		if len(term) == 1 {
			fields[t] = term
		} else {
			fields[t] = stripNonAlphaNumeric(term)
		}
	}
	return strings.Join(fields, " ")
}

func lowerCase(tokens []string) []string {
	lower := make([]string, len(tokens))
	for i, str := range tokens {
		lower[i] = strings.ToLower(str)
	}
	return lower
}

func stripNonAlphaNumeric(token string) string {
	s := []byte(token)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}
