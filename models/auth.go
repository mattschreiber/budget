package models

import (
  // "budget/middlewares"
  "strings"
  "fmt"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func VerifyLogin(user UserCredentials) (string, error) {
  var email string
  var passwordHash string
  err := db.QueryRow("SELECT email, password FROM users WHERE email = $1", strings.ToLower(user.Username)).Scan(&email, &passwordHash)
  if err != nil {
    fmt.Println(err)
    return "", err
  }
  return passwordHash, nil
}
