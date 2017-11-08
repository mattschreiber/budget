package models

import (
  "time"
  "database/sql"
  _ "github.com/lib/pq"
  "fmt"
  "sync"
)

type Budget struct {
  Id int `json:"id"`
  Credit int `json:"credit"`
  Debit int `json:"debit"`
  Trans_date time.Time `json:"trans_date"`
  Store_id int `json:"store_id"`
  User_id sql.NullInt64 `json:"user_id"`
  Category_id int `json:"category_id"`
  Applied bool `json:"applied"`
  Store_name string `json:"store_name,omitempty"`
}

func AllBudgetEntries() ([]Budget, error) {
  // now := time.Now()
  // before := time.Date(1900, 01, 15, 0, 0, 0, 0, time.UTC)
  rows, err := db.Query("SELECT b.id, b.credit, b.debit, b.trans_date, b.store_id, b.user_id, b.category_id, b.applied, s.store_name  FROM budget as b join store as s on b.store_id = s.id")
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var budgetEntries []Budget
  for rows.Next() {
    var budgetRow Budget
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_id, &budgetRow.User_id, &budgetRow.Category_id, &budgetRow.Applied, &budgetRow.Store_name)
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

// ProjectedBalance is a function that calculates the total budgeted balance at any given time.
// The function accepts a starting and ending date and returns the balance in cents as an integer
// The projected balance is calculated by first finding the current actual ledger balance.
// The second step is to calculate the budgeted balance for a given time period using a provided date range and eliminating any entries from
// the budget table where the applied column is equal to true.
// The applied column should be set to true when an entry in the budget table is recorded on the ledger.
// The last step is to sum the ledger balance with the budget balance.
func ProjectedBalance(endDate time.Time) (projBalance int, err error) {

  //sql statements that will be called concurrently to get budget and ledger balances
  budgetStmt := "SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date <= $1 AND applied != true"
  ledgerStmt := "SELECT sum(credit-debit) as balance from ledger WHERE trans_date <= $1"
  sqlStmts := []string{budgetStmt, ledgerStmt}

  var wg sync.WaitGroup
  var sumBalance int
  // anonymous function used by go routines to process ledger and budget balance queries concurrently
  for _, sqlStmt := range sqlStmts {
    wg.Add(1)
    go func(sqlStmt string, endDate time.Time){
      var balance int
      defer wg.Done()
      err = db.QueryRow(sqlStmt, endDate).Scan(&balance)
      sumBalance += balance
    }(sqlStmt, endDate)
  }
  wg.Wait()
  if err != nil {
    return -1, err
  }
  return sumBalance, nil



  // var ledgerBalance int
  // err = db.QueryRow("SELECT sum(credit-debit) as balance from ledger WHERE trans_date <= $1", endDate).Scan(&ledgerBalance)
  // if err != nil {
  //   return -1, err
  // }
  // var budgetBalance int
  // err = db.QueryRow("SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date <= $1 AND applied != true", endDate).Scan(&budgetBalance)
  // if err != nil {
  //   return -1, err
  // }
  // return ledgerBalance + budgetBalance, nil


  // type Balance struct {
  //   Amount int
  //   Error error
  // }
  // getBalance := func(sqlStmt string, endDate time.Time, c chan Balance){
  //   var balance int
  //   err = db.QueryRow(sqlStmt, endDate).Scan(&balance)
  //   if err != nil {
  //     c <- Balance{0, err}
  //   }
  //   c <- Balance{balance, nil}
  // }
  //
  // //sql statements that will be called concurrently to get budget and ledger balances
  // budgetStmt := "SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date <= $1 AND applied != true"
  // ledgerStmt := "SELECT sum(credit-debit) as balance from ledger WHERE trans_date <= $1"
  //
  //
  // c := make(chan Balance) // channel for Balance amount and error handling
  // go getBalance(budgetStmt, endDate, c)
  // go getBalance(ledgerStmt, endDate, c)
  //
  // budgetBal, ledgerBal := <-c, <-c // receive Balance struct from go routines
  //
  // if budgetBal.Error != nil {
  //   return -1, budgetBal.Error
  // }
  // if ledgerBal.Error != nil {
  //   return -1, ledgerBal.Error
  // }
  // return ledgerBal.Amount + budgetBal.Amount, nil

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
