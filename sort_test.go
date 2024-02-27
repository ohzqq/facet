package facet

var fieldSortParams = []string{}

//func TestFieldSort(t *testing.T) {
//  t.SkipNow()
//  var sorted []*txt.Token
//  alpha := libCfgStr + "&sortFacetValuesBy=alpha"
//  i, err := New(alpha)
//  if err != nil {
//    log.Fatal(err)
//  }
//  idx := NewResponse(i.Data, alpha)
//  tags := idx.GetFacet("tags")
//  for _, o := range []string{"desc", "asc"} {
//    tags.Order = o
//    sorted = tags.SortTokens()
//    switch tags.Order {
//    case "asc":
//      if sorted[0].Label != "abo" {
//        t.Errorf("alpha: %s (%d)\n", sorted[0].Label, sorted[0].Count())
//      }
//    case "desc":
//      fallthrough
//    default:
//      if sorted[0].Label != "zombies" {
//        t.Errorf("alpha: %s (%d)\n", sorted[0].Label, sorted[0].Count())
//      }
//    }
//  }

//  count := libCfgStr + "&sortFacetValuesBy=count"
//  i, err = New(count)
//  if err != nil {
//    log.Fatal(err)
//  }
//  idx = NewResponse(i.Data, alpha)
//  tags = idx.GetFacet("tags")
//  for _, o := range []string{"desc", "asc"} {
//    tags.Order = o
//    sorted = tags.SortTokens()
//    switch tags.Order {
//    case "asc":
//      if sorted[0].Label != "courting" {
//        t.Errorf("count: %s (%d)\n", sorted[0].Label, sorted[0].Count())
//      }
//    case "desc":
//      fallthrough
//    default:
//      if sorted[0].Label != "dnr" {
//        t.Errorf("count: %s (%d)\n", sorted[0].Label, sorted[0].Count())
//      }
//    }
//  }
//}
