package txt

import (
	"github.com/spf13/cast"
)

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
