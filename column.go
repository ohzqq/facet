package facet

import (
	"github.com/kelindar/column"
	"github.com/spf13/cast"
)

type Idx struct {
	*column.Collection
	pk string
}

func NewCol(pk string, cols []string, data []map[string]any, opts ...column.Options) (*Idx, error) {
	idx := &Idx{
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
