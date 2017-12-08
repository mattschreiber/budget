package models

import (
  "fmt"
  "os"
  "database/sql"
  _ "github.com/lib/pq"
)

// const (
//   host = "localhost"
//   port = 5432
//   user = "postgres"
//   dbname = "budget_development"
// )
var db *sql.DB
func InitDB() {

  host := "aatmdstad6dti5.clf0wvbr9kqy.us-east-2.rds.amazonaws.com"
  port := 5432
  user := os.Getenv("AWS_DB_USER")
  password := os.Getenv("AWS_MYBUDGET_DB_PW")
  dbname := "budget_development"

  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
    "dbname=%s sslmode=require",
    host, port, user, password, dbname)

  var err error
  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("Successfully connected!")
}
