package models

import (
  "time"
  "fmt"
)

type CategoryAmounts struct {
  Category string `json:"category"`
  Ledger int `json:"ledger"`
  Budget int `json:"budget"`
}

// This function returns a map with keys for ledger and budget. Each key contains
// an array of CategoryAmounts, which are calculated totals for a user provided time period.
func AmountsByCategory(startDate, endDate time.Time) ([]CategoryAmounts, error) {
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
      where trans_date between '2017-12-1' AND '2017-12-31'
      UNION ALL
      SELECT
        category_id, c.category_name as l_cat, 0 as actual, (credit-debit) as budget
      FROM
        budget join category c on budget.category_id = c.id
      where trans_date between $1 AND $2
  ) x
  GROUP BY
      category_id, l_cat`, startDate, endDate)

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

  // var m map[string][]CategoryAmounts
  // m = make(map[string][]CategoryAmounts)
  //
  // for _, ca := range arrCa {
  //   if ca.Budget > 0 {
  //     m["budget"] = append(m["budget"], ca)
  //   }
  //   if ca.Ledger > 0 {
  //     m["ledger"] = append(m["ledger"], ca)
  //   }
  // }
  return arrCa, nil

}
