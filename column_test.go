package facet

import "testing"

var books []map[string]any

func TestBookRow(t *testing.T) {
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
