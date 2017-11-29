package models

import (
  "time"
  "fmt"
  // "sync"
)

// type Ledger struct {
//   Id int `json:"id"`
//   Credit int `json:"credit"`
//   Debit int `json:"debit"`
//   Trans_date time.Time `json:"trans_date"`
//   Store_name string `json:"store_name"`
//   Store_id int `json:"store_id"`
// }

func AllLedgerEntries(startDate, endDate time.Time) ([]Model, error) {
  // now := time.Now()
  // before := time.Date(1900, 01, 15, 0, 0, 0, 0, time.UTC)
  rows, err := db.Query(`SELECT l.id, l.credit, l.debit, l.trans_date, s.store_name, c.category_name, l.store_id, l.category_id
    FROM ledger as l join store as s on l.store_id = s.id join category as c on l.category_id = c.id
    WHERE trans_date  BETWEEN $1 AND $2 ORDER BY l.trans_date, l.id ASC`, startDate, endDate)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var ledgerEntries []Model
  for rows.Next() {
    var ledgerRow Model
    err := rows.Scan(&ledgerRow.Id, &ledgerRow.Credit, &ledgerRow.Debit, &ledgerRow.Trans_date, &ledgerRow.St.Store_name,
      &ledgerRow.Cat.Category_name, &ledgerRow.St.Id, &ledgerRow.Cat.Id)
    if err != nil {
      return nil, err
    }
    ledgerEntries = append(ledgerEntries, ledgerRow)
  }
  if err = rows.Err(); err != nil {
   return nil, err
  }

  return ledgerEntries, nil
}
