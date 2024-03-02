package facet

import (
	"slices"
	"strings"
)

func (t *Field) SortTokens() []*Keyword {
	tokens := t.keywords

	switch t.SortBy {
	case SortByAlpha:
		if t.Order == "" {
			t.Order = "asc"
		}
		SortTokensByAlpha(tokens)
	default:
		SortTokensByCount(tokens)
	}

	if t.Order == "desc" {
		slices.Reverse(tokens)
	}

	return tokens
}

func SortTokensByCount(items []*Keyword) []*Keyword {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortTokensByAlpha(items []*Keyword) []*Keyword {
	slices.SortStableFunc(items, SortByAlphaFunc)
	return items
}

func SortByCountFunc(a *Keyword, b *Keyword) int {
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

func SortByAlphaFunc(a *Keyword, b *Keyword) int {
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
