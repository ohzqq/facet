package txt

import (
	"strings"

	"github.com/spf13/cast"
)

type Analyzer interface {
	Tokenize(any) []*Token
}

func Keyword() Analyzer {
	return keyword{}
}

//func Keyword() Option {
//  return func(tokens *Tokens) {
//    tokens.analyzer = keyword{}
//  }
//}

type keyword struct{}

func (kw keyword) Tokenize(str any) []*Token {
	return KeywordTokenizer(str)
}

func (kw keyword) Search(text string) []*Token {
	return []*Token{NewToken(normalizeText(text))}
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
