package tasks

import (
  // "fmt"
  "time"
  "budget/models"
  "budget/utilities"
)

//models.AutoPay() is found in models/ledger.go

var defaultDate = utilities.InitDates()

func ScheduledTasks() {

  t := defaultDate.Today
  n := time.Date(t.Year(), t.Month(), t.Day(), 3, 0, 0, 0, defaultDate.GetEst())

  if t.After(n)  {
    // run job immediately and then wait until 3am tomorrow
    models.AutoPay()
    n = n.Add(24 * time.Hour)
  }

  d := n.Sub(t)
  time.AfterFunc(d, tick)
}

func tick() {
  models.AutoPay()
  time.AfterFunc(time.Hour * 24, tick)
}
