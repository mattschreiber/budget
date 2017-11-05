package models

import (
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  // password = "your-password"
  dbname = "budget_development"
)
var db *sql.DB
func InitDB() {
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
  fmt.Println("Successfully connected!")
}
