package facet

import "github.com/spf13/cast"

type Facets struct {
	fields map[string]*Field
}

func New(data map[string]map[string]any, attrs []string) *Facets {
	return &Facets{
		fields: CalculateFacets(data, attrs),
	}
}

func CalculateFacets(data map[string]map[string]any, fields []string) map[string]*Field {
	facets := NewFields(fields)
	return calculateFacets(data, facets)
}

func calculateFacets(data map[string]map[string]any, facets map[string]*Field) map[string]*Field {
	for id, d := range data {
		for attr, facet := range facets {
			if val, ok := d[attr]; ok {
				facet.Add(
					val,
					[]int{cast.ToInt(id)},
				)
			}
		}
	}
	return facets
}
