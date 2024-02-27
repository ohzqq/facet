package txt

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func Fulltext() Analyzer {
	return FullText{}
}

type FullText struct{}

func (ft FullText) Tokenize(val any) []*Token {
	return FulltextTokenizer(val)
}

func (ft FullText) Search(text string) []*Token {
	return FulltextTokenizer(text)
}

func FulltextTokenizer(str any) []*Token {
	var tokens []string
	var items []*Token
	for _, token := range strings.FieldsFunc(cast.ToString(str), NotAlphaNumeric) {
		lower := strings.ToLower(token)
		if !lo.Contains(stopWords, lower) {
			items = append(items, NewToken(token))
			tokens = append(tokens, lower)
		}
	}
	for i, t := range stemmerFilter(tokens) {
		items[i].Value = t
	}
	return items
}

func rmStopWords(tokens []string) []string {
	var words []string
	for _, token := range tokens {
		if !lo.Contains(stopWords, token) {
			words = append(words, token)
		}
	}
	return words
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

func NotAlphaNumeric(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = english.Stem(token, false)
	}
	return r
}

var stopWords = []string{
	"i",
	"vol",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"of",
	"at",
	"by",
	"for",
	"with",
	"into",
	"to",
	"from",
	"then",
	"when",
	"where",
	"why",
	"how",
	"no",
	"not",
	"than",
	"too",
}
