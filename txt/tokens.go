package txt

type Tokens struct {
	tokens   map[string]*Token
	Tokens   []string
	analyzer Analyzer
}

func NewTokens() *Tokens {
	tokens := &Tokens{
		tokens:   make(map[string]*Token),
		analyzer: Simple{},
	}
	return tokens
}

func (t *Tokens) SetAnalyzer(ana Analyzer) *Tokens {
	t.analyzer = ana
	return t
}

func (t *Tokens) Find(val any) []*Token {
	var tokens []*Token
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
			t.tokens = make(map[string]*Token)
		}
		if _, ok := t.tokens[token.Value]; !ok {
			t.Tokens = append(t.Tokens, token.Label)
			t.tokens[token.Value] = token
		}
		t.tokens[token.Value].Add(ids...)
	}
}

func (t *Tokens) Tokenize(val any) []*Token {
	return t.analyzer.Tokenize(val)
}

func (t *Tokens) FindByLabel(label string) *Token {
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
