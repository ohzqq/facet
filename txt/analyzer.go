package txt

import (
	"strings"

	"github.com/spf13/cast"
)

type Analyzer interface {
	Tokenize(any) []*Token
}

type Simple struct{}

func (s Simple) Tokenize(str any) []*Token {
	return []*Token{NewToken(cast.ToString(str))}
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
