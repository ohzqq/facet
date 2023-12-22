package facet

type Facet struct {
	Name  string
	Terms []*Term
	terms map[string][]string
}

type Term struct {
	Value     any
	BelongsTo []string
}

func NewFacet(name string) *Facet {
	return &Facet{
		Name:  name,
		terms: make(map[string][]string),
	}
}

func (f *Facet) AddTerm(term string, ids ...string) *Facet {
	if _, ok := f.terms[term]; !ok {
		f.terms[term] = []string{}
	}
	f.terms[term] = append(f.terms[term], ids...)
	return f
}

func NewTerm(v any) *Term {
	return &Term{
		Value: v,
	}
}

func (t *Term) Belongs(id ...string) *Term {
	t.BelongsTo = append(t.BelongsTo, id...)
	return t
}
