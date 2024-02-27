package txt

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/sahilm/fuzzy"
)

type Token struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	bits        *roaring.Bitmap
	fuzzy.Match `json:"-"`
}

func NewToken(label string) *Token {
	return &Token{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Token) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Token) SetValue(txt string) *Token {
	f.Value = txt
	return f
}

func (f *Token) Count() int {
	return int(f.bits.GetCardinality())
}

func (f *Token) Contains(id int) bool {
	return f.bits.ContainsInt(id)
}

func (f *Token) Add(ids ...int) {
	for _, id := range ids {
		if !f.Contains(id) {
			f.bits.AddInt(id)
		}
	}
}

func (f *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		f.Label: f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}
