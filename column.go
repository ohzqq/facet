package facet

import (
	"github.com/kelindar/bitmap"
	"github.com/kelindar/column"
	"github.com/spf13/cast"
)

type Col struct {
	*column.Collection
	pk   string
	name string
}

type Row struct {
	Name string
	data bitmap.Bitmap
}

func NewRow(name string, ids []uint32) *Row {
	var bits bitmap.Bitmap
	for _, id := range ids {
		bits.Set(id)
	}

	return &Row{
		Name: name,
		data: bits,
	}
}

func NewBitmap(ids []any) bitmap.Bitmap {
	var bits bitmap.Bitmap
	for _, id := range ids {
		bits.Set(cast.ToUint32(id))
	}
	return bits
}

func NewCol(name, pk string, opts ...column.Options) *Col {
	idx := &Col{
		Collection: column.NewCollection(opts...),
		name:       name,
		pk:         pk,
	}
	return idx
}

func (c *Col) SetCols(cols []string) *Col {
	for _, col := range cols {
		c.CreateColumn(col, column.ForBool())
	}
	return c
}

func NewColz(pk string, cols []string, data []map[string]any, opts ...column.Options) (*Col, error) {
	idx := &Col{
		Collection: column.NewCollection(opts...),
		pk:         pk,
	}
	idx.CreateColumn(pk, column.ForString())
	for _, col := range cols {
		err := idx.CreateColumn(col, column.ForString())
		if err != nil {
			return nil, err
		}

		idx.Query(func(txn *column.Txn) error {
			for _, item := range data {
				if terms, ok := item[col]; ok {
					for _, term := range cast.ToStringSlice(terms) {
						idx.CreateIndex(term, col, func(r column.Reader) bool {
							return r.String() == term
						})

						_, err := idx.Insert(func(r column.Row) error {
							r.SetString(pk, cast.ToString(item[pk]))
							r.SetString(col, term)
							return nil
						})
						if err != nil {
							return err
						}
					}
				}
			}
			return nil
		})
	}
	return idx, nil
}
