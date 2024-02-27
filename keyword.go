package facet

import (
	"strings"

	"github.com/spf13/cast"
)

type Analyzer interface {
	Tokenize(any) []*Item
}

func Keyword() Analyzer {
	return keyword{}
}

type keyword struct{}

func (kw keyword) Tokenize(str any) []*Item {
	return KeywordTokenizer(str)
}

func (kw keyword) Search(text string) []*Item {
	return []*Item{NewItem(normalizeText(text))}
}

func KeywordTokenizer(val any) []*Item {
	var tokens []string
	switch v := val.(type) {
	case string:
		tokens = append(tokens, v)
	default:
		tokens = cast.ToStringSlice(v)
	}
	items := make([]*Item, len(tokens))
	for i, token := range tokens {
		items[i] = NewItem(token)
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
