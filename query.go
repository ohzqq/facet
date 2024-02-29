package facet

import "strings"

func parseQuery(query string) (string, string) {
	var attr, q string
	i := 0
	for query != "" {
		var a string
		a, query, _ = strings.Cut(query, ":")
		if a == "" {
			continue
		}
		switch i {
		case 0:
			attr = a
		case 1:
			q = a
		}
		i++
	}
	return attr, q
}
