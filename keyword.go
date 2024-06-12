package facet

import (
	"encoding/json"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"
)

type Token struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Children *Field
	bits     *roaring.Bitmap
}

func NewToken(label string) *Token {
	return &Token{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (kw *Token) Bitmap() *roaring.Bitmap {
	return kw.bits
}

func (kw *Token) SetValue(txt string) *Token {
	kw.Value = txt
	return kw
}

func (kw *Token) Items() []int {
	i := kw.bits.ToArray()
	return cast.ToIntSlice(i)
}

func (kw *Token) Count() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Len() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Contains(id int) bool {
	return kw.bits.ContainsInt(id)
}

func (kw *Token) Add(ids ...int) {
	for _, id := range ids {
		if !kw.Contains(id) {
			kw.bits.AddInt(id)
		}
	}
}

func (kw *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"count": kw.Len(),
		"value": kw.Label,
		"hits":  kw.Items(),
	}
	return json.Marshal(item)
}

func KeywordTokenizer(val any) []*Token {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Token, len(tokens))
	for i, token := range tokens {
		items[i] = NewToken(token)
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
