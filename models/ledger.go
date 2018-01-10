package models

import (
  "time"
  "fmt"
  "database/sql"
  // "sync"
)

func AllLedgerEntries(startDate, endDate time.Time) ([]Model, error) {
  // now := time.Now()
  // before := time.Date(1900, 01, 15, 0, 0, 0, 0, time.UTC)
  rows, err := db.Query(`SELECT l.id, l.credit, l.debit, l.trans_date, s.store_name, c.category_name, l.store_id, l.category_id
    FROM ledger as l join store as s on l.store_id = s.id join category as c on l.category_id = c.id
    WHERE trans_date  BETWEEN $1 AND $2 ORDER BY l.trans_date DESC, id DESC`, startDate, endDate)
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

func GetLedgerBalance(startDate time.Time, endDate time.Time, c chan Balance) {
  var balance int
  err := db.QueryRow("SELECT sum(credit-debit) as balance from ledger WHERE trans_date BETWEEN $1 AND $2", startDate, endDate).Scan(&balance)
  if err != nil {
    fmt.Println(err)
    c <- Balance{0, err}
  }
  c <- Balance{balance, nil}
}

func CreateLedgerEntry(ledger Model) (id int, err error) {
  // only care about date so set time to 0
  err = db.QueryRow("INSERT INTO ledger (credit, debit, trans_date, store_id, category_id) VALUES($1, $2, $3, $4, $5)RETURNING id",
        ledger.Credit, ledger.Debit, ledger.Trans_date.In(getEst()), ledger.St.Id, ledger.Cat.Id).Scan(&id)
  if err != nil {
    fmt.Println(err)
    return -1, err
  }
  return id, nil
}

func DeleteLedgerEntry(id string) (count int64, err error) {
  deleteEntryStmt := "DELETE FROM ledger where id = $1"
  res, err := db.Exec(deleteEntryStmt, id)
  if err != nil {
    return -1, err
  }
  count, err = res.RowsAffected()
  if err != nil {
    return -1, err
  }
  return count, nil
}

func AutoPay() {

  month := int(time.Now().Month())
  year := time.Now().Year()
  pd := time.Now()
  // firstOfMonth := time.Date(pd.Year(), pd.Month(), 1, 0, 0, 0, 0, getEst())
  // lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

  // find all budgetEntries with a trans_date of today
  budgetEntries, err := AutoPayBudgetEntries(pd)
  if err != nil {
    fmt.Println(err)
    return
  }
  // if there are entries for today
  if budgetEntries != nil {
    //loop through todays entries, check if a corresponding entry exists on the ledger for current monthly
    // and if no ledger entry, then create one
    for _, entry := range budgetEntries {
      var ledgerEntry Model
      err = db.QueryRow(`SELECT id, store_id FROM ledger
        WHERE extract(month from trans_date) = $1 AND extract(year from trans_date) = $2
        AND ledger.store_id = $3`, month, year, entry.St.Id).Scan(&ledgerEntry.Id, &ledgerEntry.St.Id)
      if err != nil {
        if err == sql.ErrNoRows {
          insertEntryStmt := "INSERT INTO ledger (credit, debit, trans_date, store_id, category_id) VALUES ($1, $2, $3, $4, $5)"
          res, err := db.Exec(insertEntryStmt, entry.Credit, entry.Debit, time.Now(), entry.St.Id, entry.Cat.Id)
          if err != nil {
            fmt.Println(err)
          }
          _, err = res.RowsAffected()
          if err != nil {
            fmt.Println(err)
            return
          }
          // This is good place to send build email to send with info on new ledger entries
        } else {
          fmt.Println(err)
        }
      }
    }
  }
}
