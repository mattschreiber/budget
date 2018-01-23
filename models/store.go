package models

import (
  "fmt"
)

type Store struct {
  Id int `json:"id"`
  Store_name string `json:"store_name"`
  Default_credit int `json:"default_credit"`
  Default_debit int `json:"default_debit"`
}

// a function that returns a slice containing all stores from the store table
func GetAllStores() ([]Store, error) {
  rows, err := db.Query("SELECT s.id, s.store_name, s.default_credit, s.default_debit FROM store as s order by s.store_name")
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var stores []Store
  for rows.Next() {
    var store Store
    err := rows.Scan(&store.Id, &store.Store_name, &store.Default_credit, &store.Default_debit)
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

func CreateStore(store Store) (id int, err error) {
  // only care about date so set time to 0
  err = db.QueryRow("INSERT INTO store (store_name) VALUES($1)RETURNING id",
        store.Store_name).Scan(&id)
  if err != nil {
    fmt.Println(err)
    return -1, err
  }
  return id, nil
}
