package facet

import (
	"slices"
	"strings"
)

func (f *Field) SortTokens() []*Token {
	tokens := f.keywords

	switch f.SortBy {
	case SortByAlpha:
		if f.Order == "" {
			f.Order = "asc"
		}
		SortTokensByAlpha(tokens)
	default:
		SortTokensByCount(tokens)
	}

	if f.Order == "desc" {
		slices.Reverse(tokens)
	}

	return tokens
}

func SortTokensByCount(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortTokensByAlpha(items []*Token) []*Token {
	slices.SortStableFunc(items, SortByAlphaFunc)
	return items
}

func SortByCountFunc(a *Token, b *Token) int {
	aC := a.Count()
	bC := b.Count()
	switch {
	case aC > bC:
		return 1
	case aC == bC:
		return 0
	default:
		return -1
	}
}

func SortByAlphaFunc(a *Token, b *Token) int {
	aL := strings.ToLower(a.Label)
	bL := strings.ToLower(b.Label)
	switch {
	case aL > bL:
		return 1
	case aL == bL:
		return 0
	default:
		return -1
	}
}
