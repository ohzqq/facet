package facet

import (
	"testing"

	"github.com/kelindar/bitmap"
)

var books []map[string]any

const numBooks = 7174

func TestIdxBitmap(t *testing.T) {
	ids := idx.Bitmap()
	if ids.Count() != len(idx.items) {
		t.Errorf("got %d rows, expected %d\n", ids.Count(), len(idx.items))
	}
}

func TestOrFilter(t *testing.T) {
	abo := aboFilter(t)
	dnr := dnrFilter(t)
	abo.Or(dnr)
	if abo.Count() != 2269 {
		t.Errorf("got %d, expected %d\n", abo.Count(), 2269)
	}
}

func TestAndFilter(t *testing.T) {
	ids := aboFilter(t)
	dnr := dnrFilter(t)
	ids.And(dnr)
	if ids.Count() != 384 {
		t.Errorf("got %d, expected %d\n", ids.Count(), 384)
	}
}

func TestDnrFilter(t *testing.T) {
	ids := dnrFilter(t)
	if ids.Count() != 2237 {
		t.Errorf("got %d, expected %d\n", ids.Count(), 2237)
	}
}

func TestAboFilter(t *testing.T) {
	ids := aboFilter(t)
	if ids.Count() != 416 {
		t.Errorf("got %d, expected %d\n", ids.Count(), 416)
	}
}

func aboFilter(t *testing.T) bitmap.Bitmap {
	ids := idx.Bitmap()
	term := idx.GetTerm("tags", "abo")
	bits := term.Bitmap()
	if term.Count != bits.Count() {
		t.Errorf("got %d items, expected %d\n", bits.Count(), term.Count)
	}
	ids.Filter(func(x uint32) bool {
		return bits.Contains(x)
	})
	return ids
}

func dnrFilter(t *testing.T) bitmap.Bitmap {
	ids := idx.Bitmap()
	term := idx.GetTerm("tags", "dnr")
	bits := term.Bitmap()
	if term.Count != bits.Count() {
		t.Errorf("got %d items, expected %d\n", bits.Count(), term.Count)
	}
	ids.Filter(func(x uint32) bool {
		return bits.Contains(x)
	})
	return ids
}

//func TestNewCol(t *testing.T) {
//  books := loadData(t)
//  //col, err := NewColz("id", []string{"authors", "tags"}, books)
//  //if err != nil {
//  //t.Error(err)
//  //}

//  tags := NewCol("tags", "id")
//  terms := CollectTerms("tags", books)
//  //println(len(terms))

//  tags.SetCols(terms)

//  fmt.Println(tags.Count())

//  //col.Query(func(txn *column.Txn) error {
//  //count := txn.WithValue("tags", func(v interface{}) bool {
//  //return v == "abo"
//  //}).Count()
//  //println(count)
//  //return nil
//  //})

//  // How many rogues and mages?
//  //col.Query(func(txn *column.Txn) error {
//  //c := txn.With("abo").Union("dnr").Count()
//  //println(c)
//  //return nil
//  //})

//  err := tags.Close()
//  if err != nil {
//    t.Error(err)
//  }
//}
