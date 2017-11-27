package models

import (
  "time"
  "database/sql"
  _ "github.com/lib/pq"
  "fmt"
  // "sync"
)

type Budget struct {
  Id int `json:"id"`
  Credit int `json:"credit"`
  Debit int `json:"debit"`
  Trans_date time.Time `json:"trans_date"`
  Store_id int `json:"store_id"`
  User_id sql.NullInt64 `json:"user_id"`
  Category_id int `json:"category_id"`
  Store_name string `json:"store_name,omitempty"`
}

type Balance struct {
  Amount int
  Error error
}

func AllBudgetEntries(startDate, endDate time.Time) ([]Budget, error) {
  // now := time.Now()
  // before := time.Date(1900, 01, 15, 0, 0, 0, 0, time.UTC)
  rows, err := db.Query(`SELECT b.id, b.credit, b.debit, b.trans_date, b.store_id, b.user_id, b.category_id, s.store_name FROM budget as b join store as s on b.store_id = s.id WHERE trans_date BETWEEN $1 AND $2`, startDate, endDate)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var budgetEntries []Budget
  for rows.Next() {
    var budgetRow Budget
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_id, &budgetRow.User_id, &budgetRow.Category_id, &budgetRow.Store_name)
    if err != nil {
      fmt.Println("error scaning", err)
      return nil, err
    }
    budgetEntries = append(budgetEntries, budgetRow)
  }
  if err = rows.Err(); err != nil {
    fmt.Println(err)
   return nil, err
  }

  return budgetEntries, nil
}

func BudgetEntry(id int) (budget Budget, err error) {
  err = db.QueryRow("SELECT * FROM budget WHERE id = $1", id).Scan(&budget.Id, &budget.Credit, &budget.Debit, &budget.Trans_date, &budget.Store_id, &budget.User_id, &budget.Category_id)
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

func GetBudgetBalance(startDate time.Time, endDate time.Time, c chan Balance) {
  var balance int
  err := db.QueryRow("SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date BETWEEN $1 AND $2",
        startDate, endDate).Scan(&balance)
  if err != nil {
    fmt.Println(err)
    c <- Balance{0, err}
  }
  c <- Balance{balance, nil}
}

func GetLedgerBalance(startDate time.Time, endDate time.Time, c chan Balance) {
  var balance int
  err := db.QueryRow("SELECT sum(credit-debit) as balance from ledger WHERE trans_date > $1 AND trans_date < $2", startDate, endDate).Scan(&balance)
  if err != nil {
    fmt.Println(err)
    c <- Balance{0, err}
  }
  c <- Balance{balance, nil}
}
// ProjectedBalance is a function that calculates the total budgeted balance at any given time.
// The function accepts an ending date and returns the balance in cents as an integer
// The projected balance is calculated by first finding the current actual ledger balance up through the end of the
// previous pay period.
// The second step is to calculate the budgeted balance starting with the current pay period and up until a user provided end Date
// The last step is to sum the ledger balance with the budget balance.
func ProjectedBalance(endDate time.Time) (projBalance int, err error) {
  c := make(chan Balance) // channel for Balance amount and error handling

  // If the end date provided is after the current pay period then find all budget entries starting with the
  // the beginning of the next pay period and ending on the user provided end date
  if endDate.After(currentPayPeriod(time.Now())) {
    go GetBudgetBalance(currentPayPeriod(time.Now()), endDate, c)
  }else {
    go GetBudgetBalance(prevPayDate(time.Now()), endDate, c)
  }
  // If the end date for the projection time period is before the current pay period ends then
  // only include ledger entries through the previous pay period. Otherwise include all entries
  // until the provided end date. The first parament is just a dummy date used to get the entire ledger balance.
  if endDate.Before(currentPayPeriod(time.Now())) {
    go GetLedgerBalance(beginningOfTime(), prevPayDate(time.Now()), c)
  }else {
    go GetLedgerBalance(beginningOfTime(), endDate, c)
  }
  budgetBal, ledgerBal := <-c, <-c // receive Balance struct from go routines
  //
  if budgetBal.Error != nil || ledgerBal.Error != nil {
    if budgetBal.Error != nil {
      return -1, budgetBal.Error
    }else {
      return -1, ledgerBal.Error
    }
  }

  return ledgerBal.Amount + budgetBal.Amount, nil
}

type TotalAmounts struct {
  LedgerAmount int `json:"ledgeramount"`
  BudgetAmount int `json:"budgetamount"`
}
// This is a function that will return both the ledger and budget totals for a given time period
func GetSpending(startDate time.Time, endDate time.Time) (total TotalAmounts, err error) {

  c1 := make(chan Balance) // channel for Balance amount and error handling
  c2 := make(chan Balance) // channel for Balance amount and error handling

  go GetBudgetBalance(startDate, endDate, c1)
  go GetLedgerBalance(startDate, endDate, c2)

  for i := 0; i < 2; i++ {
    select {
    case budgetBal := <-c1:
      total.BudgetAmount = budgetBal.Amount
    case ledgerBal := <-c2:
      total.LedgerAmount = ledgerBal.Amount
    }
  }
  return total, nil
}
// utility function that finds the end of the previous pay period which is needed to calculate the future projected balance
// pay periods are assumed to be bi-monthly
func prevPayDate(today time.Time) (time.Time) {
  // today := time.Now()
  var prevPayDate time.Time
  if today.Day() >= 15 {
    middleOfMonth := time.Date(today.Year(), today.Month(), 15, 0, 0, 0, 0, time.UTC)
    prevPayDate = middleOfMonth
  }else {
    currentYear, currentMonth, _ := today.Date()
    currentLocation := today.Location()
    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    // lastOfMonth := firstOfMonth.AddDate(0, 0, 0)
    prevPayDate = firstOfMonth
  }
  return prevPayDate
}

func currentPayPeriod(today time.Time) (time.Time) {

  var endPayPeriod time.Time
  if today.Day() >= 15 {
    // endPayPeriod is equal to first day of next month if date >= 15
    endPayPeriod = time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, time.UTC)
  }else {
    // endPayPeriod equal to 15th of the month.
    endPayPeriod = time.Date(today.Year(), today.Month(), 15, 0, 0, 0, 0, time.UTC)
  }
  return endPayPeriod
}

//This is used when a default date is needed for the start of a date range
func beginningOfTime() (time.Time) {
  beginningOfTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
  return beginningOfTime
}

  // //sql statements that will be called concurrently to get budget and ledger balances
  // budgetStmt := "SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date >= $1"
  // ledgerStmt := "SELECT sum(credit-debit) as balance from ledger WHERE trans_date <= $1"
  // sqlStmts := []string{budgetStmt, ledgerStmt}
  //
  // var wg sync.WaitGroup
  // var sumBalance int
  // // anonymous function used by go routines to process ledger and budget balance queries concurrently
  // for _, sqlStmt := range sqlStmts {
  //   wg.Add(1)
  //   go func(sqlStmt string, endDate time.Time){
  //     var balance int
  //     defer wg.Done()
  //     err = db.QueryRow(sqlStmt, endDate).Scan(&balance)
  //     sumBalance += balance
  //   }(sqlStmt, endDate)
  // }
  // wg.Wait()
  // if err != nil {
  //   return -1, err
  // }
  // return sumBalance, nil




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


  // func ProjectedBalance(endDate time.Time) (projBalance int, err error) {
  //
  //   type Balance struct {
  //     Amount int
  //     Error error
  //   }
  //   //sql statements that will be called concurrently to get ledger balances
  //   // we search for all ledger entries up through the last pay period. we do not want current  pay period
  //   // ledger transaction when calculating the projected balance.
  //   ledgerStmt := "SELECT sum(credit-debit) as balance from ledger WHERE trans_date < $1"
  //   getLedgerBalance := func(sqlStmt string, prevPayDate time.Time, endDate time.Time, c chan Balance){
  //     var balance int
  //
  //     // If the end date for the projection time period is before the current pay period ends then
  //     // only include ledger entries through the previous pay period. Otherwise include all entries
  //     // until the provided end date.
  //     if endDate.Before(currentPayPeriod()) {
  //       err = db.QueryRow(sqlStmt, prevPayDate).Scan(&balance)
  //     }else {
  //       err = db.QueryRow(sqlStmt, endDate).Scan(&balance)
  //     }
  //     if err != nil {
  //       c <- Balance{0, err}
  //     }
  //     c <- Balance{balance, nil}
  //   }
  //   //sql statements that will be called concurrently to get budget balances
  //   // we want transactions starting with the beginning of the current pay period and up until a user provided end date
  //   budgetStmt := "SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date BETWEEN $1 AND $2"
  //   getBudgetBalance := func(sqlStmt string, prevPayDate time.Time, endDate time.Time, c chan Balance){
  //     var balance int
  //
  //     // If the end date provided is after the current pay period then find all budget entries starting with the
  //     // the beginning of the next pay period and ending on the user provided end date
  //     if endDate.After(currentPayPeriod()) {
  //       err = db.QueryRow(sqlStmt, currentPayPeriod(), endDate).Scan(&balance)
  //     }else {
  //       err = db.QueryRow(sqlStmt, prevPayDate, endDate).Scan(&balance)
  //     }
  //     if err != nil {
  //       c <- Balance{0, err}
  //     }
  //     c <- Balance{balance, nil}
  //   }
  //
  //   c := make(chan Balance) // channel for Balance amount and error handling
  //   go getBudgetBalance(budgetStmt, prevPayDate(), endDate, c)
  //   go getLedgerBalance(ledgerStmt, prevPayDate(), endDate, c)
  //   //
  //   budgetBal, ledgerBal := <-c, <-c // receive Balance struct from go routines
  //   //
  //   if budgetBal.Error != nil {
  //     return -1, budgetBal.Error
  //   }
  //   if ledgerBal.Error != nil {
  //     return -1, ledgerBal.Error
  //   }
  //   return ledgerBal.Amount + budgetBal.Amount, nil
  //
  // }
