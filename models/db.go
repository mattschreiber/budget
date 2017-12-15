package models

import (
  "fmt"
  // "os"
  "database/sql"
  _ "github.com/lib/pq"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  dbname = "budget_development"
)
var db *sql.DB
func InitDB() {

  // host := "aa1bw5g82fgor1s.clf0wvbr9kqy.us-east-2.rds.amazonaws.com"
  // port := 5432
  // user := os.Getenv("AWS_DB_USER")
  // password := os.Getenv("AWS_MYBUDGET_DB_PW")
  // dbname := "budget_development"

  // psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
  //   "dbname=%s sslmode=require",
  //   host, port, user, password, dbname)

  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
    host, port, user, dbname)

  var err error
  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  db.SetMaxIdleConns(0)
  fmt.Println("Successfully connected!")
}
