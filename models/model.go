package models

import (
  "time"
)
// Generic model for describing a budget or ledger
type Model struct {
  Id int `json:"id"`
  Credit int `json:"credit"`
  Debit int `json:"debit"`
  Trans_date time.Time `json:"trans_date"`
  Cat Category `json:"category"`
  St Store `json:"store"`
}

type Balance struct {
  Amount int
  Error error
}

type TotalAmounts struct {
  LedgerAmount int `json:"ledgeramount"`
  BudgetAmount int `json:"budgetamount"`
}

type StoreCat struct {
  Cat []Category `json:"category"`
  St []Store `json:"store"`
}
// This is a function that will return both the ledger and budget totals for a given time period
func GetAmountSpent(startDate time.Time, endDate time.Time) (total TotalAmounts, err error) {

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
    // adjust date for ledger entry in order to not count pay dates twice.
    // This is necessary in order to reuse the GetLedgerBalance and getBudgetBalance
    // The sql stmt used is BETWEEN which is why dates need to be adjusted
    pd := time.Now()
    pd = time.Date(pd.Year(), pd.Month(), 14, 0, 0, 0, 0, time.UTC) // ledger range should end on 14th of this month
    firstOfMonth := time.Date(pd.Year(), pd.Month(), 1, 0, 0, 0, 0, time.UTC)
    lastOfPrevMonth := firstOfMonth.AddDate(0, 0, -1) // ledger range should be until end of last month
    if currentPayPeriod(time.Now()).Day() == 1 {
      go GetLedgerBalance(beginningOfTime(), pd, c)
    } else {
      go GetLedgerBalance(beginningOfTime(), lastOfPrevMonth, c)
    }
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

// function that returns an object consisting of all categories and all stores
func GetStoreCategory() (StoreCat, error){
  // c1 := make(chan []Category) // channel for Balance amount and error handling
  // c2 := make(chan []Store) // channel for Balance amount and error handling
  var storeCat StoreCat
  var err error
  storeCat.Cat, err =  GetAllCategories()
  if err != nil {
    return storeCat, err
  }
  var err1 error
  storeCat.St, err1 = GetAllStores()
  if err1 != nil {
    return storeCat, err1
  }
  // go GetAllStores(c2)

  // for i := 0; i < 2; i++ {
  //   select {
  //   case categories := <-c1:
  //     StoreCat.Cat = categories
  //   case stores := <-c2:
  //     StoreCat.St = stores
  //   }
  // }
  return storeCat, nil
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
