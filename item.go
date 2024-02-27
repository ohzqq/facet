package facet

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring"
)

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
	bits  *roaring.Bitmap
}

func NewItem(label string) *Item {
	return &Item{
		Value: label,
		Label: label,
		bits:  roaring.New(),
	}
}

func (f *Item) Bitmap() *roaring.Bitmap {
	return f.bits
}

func (f *Item) SetValue(txt string) *Item {
	f.Value = txt
	return f
}

func (f *Item) Count() int {
	return int(f.bits.GetCardinality())
}

func (f *Item) Contains(id int) bool {
	return f.bits.ContainsInt(id)
}

func (f *Item) Add(ids ...int) {
	for _, id := range ids {
		if !f.Contains(id) {
			f.bits.AddInt(id)
		}
	}
}

func (f *Item) MarshalJSON() ([]byte, error) {
	item := map[string]any{
		f.Label: f.Count(),
	}
	d, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return d, nil
}
