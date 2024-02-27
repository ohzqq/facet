package facet

import (
	"slices"
	"strings"
)

func (t *Field) SortTokens() []*Item {
	tokens := t.GetTokens()

	switch t.SortBy {
	case SortByAlpha:
		if t.Order == "" {
			t.Order = "asc"
		}
		SortTokensByAlpha(tokens)
	default:
		if t.Order == "" {
			t.Order = "desc"
		}
		SortTokensByCount(tokens)
	}

	if t.Order == "desc" {
		slices.Reverse(tokens)
	}

	return tokens
}

func SortTokensByCount(items []*Item) []*Item {
	slices.SortStableFunc(items, SortByCountFunc)
	return items
}

func SortTokensByAlpha(items []*Item) []*Item {
	slices.SortStableFunc(items, SortByAlphaFunc)
	return items
}

func SortByCountFunc(a *Item, b *Item) int {
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

func SortByAlphaFunc(a *Item, b *Item) int {
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
