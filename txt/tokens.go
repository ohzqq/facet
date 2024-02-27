package txt

type Tokens struct {
	tokens   map[string]*Item
	Tokens   []string
	analyzer Analyzer
}

func NewTokens() *Tokens {
	tokens := &Tokens{
		tokens:   make(map[string]*Item),
		analyzer: Keyword(),
	}
	return tokens
}

func (t *Tokens) SetAnalyzer(ana Analyzer) *Tokens {
	t.analyzer = ana
	return t
}

func (t *Tokens) Find(val any) []*Item {
	var tokens []*Item
	for _, tok := range t.Tokenize(val) {
		if token, ok := t.tokens[tok.Value]; ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func (t *Tokens) Add(val any, ids []int) {
	for _, token := range t.Tokenize(val) {
		if t.tokens == nil {
			t.tokens = make(map[string]*Item)
		}
		if _, ok := t.tokens[token.Value]; !ok {
			t.Tokens = append(t.Tokens, token.Label)
			t.tokens[token.Value] = token
		}
		t.tokens[token.Value].Add(ids...)
	}
}

func (t *Tokens) Tokenize(val any) []*Item {
	return t.analyzer.Tokenize(val)
}

func (t *Tokens) FindByLabel(label string) *Item {
	for _, token := range t.tokens {
		if token.Label == label {
			return token
		}
	}
	return NewToken(label)
}

func (t *Tokens) Count() int {
	return len(t.tokens)
}
