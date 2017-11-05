package main

import (
  "budget/models"
)

func main() {
  models.InitDB()
  models.AllBudgetEntries()
}
