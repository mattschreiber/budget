package models

import (
  "time"
  "database/sql"
  _ "github.com/lib/pq"
)

type Budget struct {
  Id int
  Credit int
  Debit int
  Store_id int
  Category_id int
  Applied bool
  User_id sql.NullInt64
  Trans_date time.Time
}

func AllBudgetEntries () ([]Budget, error) {
  rows, err := db.Query("SELECT * FROM budget")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var budgetEntries []Budget
  for rows.Next() {
    var budgetRow Budget
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_id, &budgetRow.User_id, &budgetRow.Category_id, &budgetRow.Applied)
    if err != nil {
      return nil, err
    }
    budgetEntries = append(budgetEntries, budgetRow)
  }
  if err = rows.Err(); err != nil {
   return nil, err
  }

  return budgetEntries, nil
}

func BudgetTotal() (balance int, err error) {
  rows, err := db.Query("SELECT SUM(credit - debit) as balance FROM budget")
  if err != nil {
    return -1, err
  }
  defer rows.Close()
  for rows.Next(){
    err = rows.Scan(&balance)
    if err != nil {
      return -1, err
    }
  }
  if err = rows.Err(); err != nil {
   return -1, err
  }

  return balance, nil
}
