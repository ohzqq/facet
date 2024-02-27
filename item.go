package facet

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
)

type Keyword struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Children *Field
	bits     *roaring.Bitmap
}

func NewItem(label string) *Keyword {
	return &Keyword{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Keyword) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Keyword) SetValue(txt string) *Keyword {
	f.Value = txt
	return f
}

func (f *Keyword) Count() int {
	return int(f.bits.GetCardinality())
}

func (f *Keyword) Contains(id int) bool {
	return f.bits.ContainsInt(id)
}

func (f *Keyword) Add(ids ...int) {
	for _, id := range ids {
		if !f.Contains(id) {
			f.bits.AddInt(id)
		}
	}
}

func (f *Keyword) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		f.Label: f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}
