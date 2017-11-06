package models

import (
  "time"
  "database/sql"
  _ "github.com/lib/pq"
  "fmt"
)

type Budget struct {
  Id int
  Credit int
  Debit int
  Trans_date time.Time
  Store_id int
  User_id sql.NullInt64
  Category_id int
  Applied bool
}

func AllBudgetEntries() ([]Budget, error) {
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

func BudgetEntry(id int) (budget Budget, err error) {
  err = db.QueryRow("SELECT * FROM budget WHERE id = $1", id).Scan(&budget.Id, &budget.Credit, &budget.Debit, &budget.Trans_date, &budget.Store_id, &budget.User_id, &budget.Category_id, &budget.Applied)
  if err != nil {
    fmt.Println(err)
    return Budget{}, err
  }
  return budget, nil
}
func BudgetTotal(t time.Time) (balance int, err error) {

  rows, err := db.Query("SELECT SUM(credit - debit) as balance FROM budget where trans_date <= $1", t)
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

func ProjectedBalance() (projBalance int, err error) {
  err = db.QueryRow("SELECT currentBudget();").Scan(&projBalance)
  if err != nil {
    return -1, err
  }
  return projBalance, nil
}

// if t.Day() < 15 {
//   middleOfMonth := time.Date(t.Year(), t.Month(), 15, 0, 0, 0, 0, time.UTC)
//   fmt.Println(middleOfMonth)
// }else {
//   currentYear, currentMonth, _ := t.Date()
//   currentLocation := t.Location()
//   firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
//   lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
//   fmt.Println(lastOfMonth)
// }
