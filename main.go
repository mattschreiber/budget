package main

import (
	"database/sql"
	"fmt"
  "log"
  "time"
	_ "github.com/lib/pq"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  // password = "your-password"
  dbname = "budget_development"
)

type Budget struct {
  Id int
  Credit int
  Debit int
  Store_id int
  User_id sql.NullInt64
  Trans_date time.Time
}

func main() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
    host, port, user, dbname)

  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("Successfully connected!")

  rows, err := db.Query("SELECT * FROM budget")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  var budgetEntries []Budget
  for rows.Next() {
    var budgetRow Budget
    err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_id, &budgetRow.User_id)
    if err != nil {
      log.Fatal(err)
    }
    budgetEntries = append(budgetEntries, budgetRow)
  }
  if err = rows.Err(); err != nil {
   log.Fatal(err)
  }
  for _, v := range budgetEntries {
    fmt.Printf("%d, %d, %d, %d, %s\n", v.Id, v.Credit, v.Debit, v.Store_id, v.Trans_date.Format(time.RFC822))
  }
}
