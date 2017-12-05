package models

import (
  "time"
  _ "github.com/lib/pq"
  "fmt"
  // "sync"
)

func AllBudgetEntries(startDate, endDate time.Time) ([]Model, error) {
  // now := time.Now()
  rows, err := db.Query(`SELECT b.id, b.credit, b.debit, b.trans_date, s.store_name, c.category_name, b.store_id, b.category_id
    FROM budget as b join store as s on b.store_id = s.id join category as c on b.category_id = c.id
    WHERE b.trans_date  BETWEEN $1 AND $2 ORDER BY b.trans_date DESC, id DESC`, startDate, endDate)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var budgetEntries []Model
  for rows.Next() {
    var budgetRow Model
    // err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_name, &budgetRow.Category_name, &budgetRow.Store_id, &budgetRow.Category_id)
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.St.Store_name, &budgetRow.Cat.Category_name, &budgetRow.St.Id, &budgetRow.Cat.Id)
    if err != nil {
      fmt.Println("error scaninng", err)
      return nil, err
    }
    budgetEntries = append(budgetEntries, budgetRow)
  }
  if err = rows.Err(); err != nil {
    fmt.Println("error scanning a row", err)
   return nil, err
  }
  return budgetEntries, nil
}

// func BudgetEntry(id int) (budget Budget, err error) {
//   err = db.QueryRow("SELECT * FROM budget WHERE id = $1", id).Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_name, &budgetRow.Category_name, &budgetRow.Store_id, &budgetRow.Category_id)
//   if err != nil {
//     fmt.Println(err)
//     return Budget{}, err
//   }
//   return budget, nil
// }
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

func CreateBudgetEntry(budget Model) (id int, err error) {
  // only care about date so set time to 0
  err = db.QueryRow("INSERT INTO budget (credit, debit, trans_date, store_id, category_id) VALUES($1, $2, $3, $4, $5)RETURNING id",
        budget.Credit, budget.Debit, budget.Trans_date.In(getEst()), budget.St.Id, budget.Cat.Id).Scan(&id)
  if err != nil {
    fmt.Println(err)
    return -1, err
  }
  return id, nil
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
