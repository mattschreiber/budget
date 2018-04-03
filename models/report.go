package models

import (
  "fmt"
  "time"
)

type CategoryAmounts struct {
  Category string `json:"category"`
  Ledger int `json:"ledger"`
  Budget int `json:"budget"`
}

// This function returns a map with keys for ledger and budget. Each key contains
// an array of CategoryAmounts, which are calculated totals for a user provided time period.
func AmountsByCategory(month, year string) ([]CategoryAmounts, error) {
  rows, err := db.Query(`SELECT
    l_cat,
    ABS(sum(actual)) as actual,
    ABS(sum(budget)) as budget
  FROM
  (
      SELECT
        category_id, c.category_name as l_cat, (credit-debit) as actual , 0 as budget
      FROM
        ledger join category c on ledger.category_id = c.id
      where extract(month from trans_date) = $1 AND extract(year from trans_date) = $2
      UNION ALL
      SELECT
        category_id, c.category_name as l_cat, 0 as actual, (credit-debit) as budget
      FROM
        budget join category c on budget.category_id = c.id
        where extract(month from trans_date) = $1 AND extract(year from trans_date) = $2
  ) x
  GROUP BY
      category_id, l_cat`, month, year)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var arrCa []CategoryAmounts

  for rows.Next() {
    var ca CategoryAmounts
    err := rows.Scan(&ca.Category, &ca.Ledger, &ca.Budget)
    if err != nil {
      fmt.Println("err: ", err)
      return nil, err
    }
    arrCa = append(arrCa, ca)
  }

  if err = rows.Err(); err != nil {
    fmt.Println("error scanning a row", err)
   return nil, err
  }

  return arrCa, nil

}

type MonthTotals struct {
  BudgetTotal int `json:"budget_total"`
  LedgerTotal int `json:"ledger_total"`
  Month string `json:"month"`
  Year float64 `json:"year"`
  // NumMonth int `json:"num_month"`
}

// returns monthly amounts for budget and ledger. calculated by subtracting total debits from total credits
// query assumes >= startDate and < endDate
func MonthlyTotalSpent(startDate, endDate time.Time) ([]MonthTotals, error) {
  // query to return budget and ledger amount spent by month and year
  rows, err := db.Query(`SELECT
      mon,
      yyyy,
      sum(actual) as actual,
      sum(budget) as budget
  FROM
  (
      SELECT
        to_char(trans_date,'Mon') as mon, extract(year from trans_date) as yyyy,
        sum(credit-debit) as actual , 0 as budget, date_trunc('month', trans_date) as num_month
      FROM ledger
      WHERE trans_date >= $1 AND trans_date < $2
      GROUP BY mon, yyyy, date_trunc('month', trans_date)
      UNION ALL
      SELECT to_char(trans_date,'Mon') as mon, extract(year from trans_date) as yyyy,
        0 as actual, sum(credit-debit) as budget, date_trunc('month', trans_date) as num_month
      FROM budget
      WHERE trans_date >= $1 AND trans_date < $2
      GROUP BY mon, yyyy, date_trunc('month', trans_date)
  ) x
  GROUP BY mon, yyyy, num_month
  ORDER BY yyyy, num_month ASC`, startDate, endDate)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer rows.Close()

  var arrMonthTotals []MonthTotals

  for rows.Next() {
    var mt MonthTotals
    err := rows.Scan(&mt.Month, &mt.Year, &mt.LedgerTotal, &mt.BudgetTotal)
    if err != nil {
      fmt.Println("err: ", err)
      return nil, err
    }
    arrMonthTotals = append(arrMonthTotals, mt)
  }

  if err = rows.Err(); err != nil {
    fmt.Println("error scanning a row", err)
   return nil, err
  }

  return arrMonthTotals, nil
}
