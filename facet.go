package facet

type Facet struct {
	Name     string
	Elements []any
}

func NewFacet(name string, el []any) *Facet {
	return &Facet{
		Name:     name,
		Elements: el,
	}
}
