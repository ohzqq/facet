//go:build exclude

package facet

import (
	"fmt"
	"testing"

	"github.com/kelindar/column"
)

var col *Idx

func TestNewCol(t *testing.T) {
	books := loadData(t)
	col, err := NewCol("id", []string{"authors", "tags"}, books)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(col.Count())

	col.Query(func(txn *column.Txn) error {
		count := txn.WithValue("tags", func(v interface{}) bool {
			return v == "abo"
		}).Count()
		println(count)
		return nil
	})

	// How many rogues and mages?
	col.Query(func(txn *column.Txn) error {
		c := txn.With("abo").Union("dnr").Count()
		println(c)
		return nil
	})

	err = col.Close()
	if err != nil {
		t.Error(err)
	}
}
