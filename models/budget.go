package models

import (
  "fmt"
  "time"
  "database/sql"
  _ "github.com/lib/pq"
  "log"
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

func AllBudgetEntries () {
  rows, err := db.Query("SELECT * FROM budget")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  var budgetEntries []Budget
  for rows.Next() {
    var budgetRow Budget
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_id, &budgetRow.User_id, &budgetRow.Category_id, &budgetRow.Applied)
    if err != nil {
      log.Fatal(err)
    }
    budgetEntries = append(budgetEntries, budgetRow)
  }
  if err = rows.Err(); err != nil {
   log.Fatal(err)
  }
  for _, v := range budgetEntries {
    fmt.Printf("%d, %d, %d, %d, %d, %t, %s\n", v.Id, v.Credit, v.Debit, v.Store_id, v.Category_id, v.Applied, v.Trans_date.Format(time.RFC822))
  }
}
