package models

import (
  "fmt"
)

type Store struct {
  Id int `json:"id"`
  Store_name string `json:"store_name"`
}

// a function that returns a slice containing all stores from the store table
func GetAllStores() ([]Store, error) {
  rows, err := db.Query("SELECT s.id, s.store_name FROM store as s order by s.store_name")
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var stores []Store
  for rows.Next() {
    var store Store
    err := rows.Scan(&store.Id, &store.Store_name)
    if err != nil {
      return nil, err
    }
    stores = append(stores, store)
  }
  if err = rows.Err(); err != nil {
   return nil, err
  }
  return stores, nil
}
