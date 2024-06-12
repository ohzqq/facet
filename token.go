package facet

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
	"github.com/avito-tech/normalize"
	"github.com/spf13/cast"
)

type Token struct {
	Value string `json:"value"`
	Label string `json:"label"`
	bits  *roaring.Bitmap
}

func NewToken(label string) *Token {
	tok := &Token{
		Label: label,
		bits:  roaring.New(),
	}
	tok.Value = normalize.Normalize(label)
	return tok
}

func (kw *Token) Bitmap() *roaring.Bitmap {
	return kw.bits
}

func (kw *Token) SetValue(txt string) *Token {
	kw.Value = txt
	return kw
}

func (kw *Token) Items() []int {
	return cast.ToIntSlice(kw.bits.ToArray())
}

func (kw *Token) Count() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Len() int {
	return int(kw.bits.GetCardinality())
}

func (kw *Token) Add(ids ...int) {
	for _, id := range ids {
		if !kw.bits.ContainsInt(id) {
			kw.bits.AddInt(id)
		}
	}
}

func (kw *Token) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		"count": kw.Len(),
		"value": kw.Label,
		"hits":  kw.Items(),
	}
	return json.Marshal(item)
}
