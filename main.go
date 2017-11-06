package main

import (
  "encoding/json"
  "budget/models"
  "fmt"
  "time"
  "os"
)

func main() {
  models.InitDB()

  budgetEntries, _ := models.AllBudgetEntries()
  for _, v := range budgetEntries {
    fmt.Printf("%d, %d, %d, %d, %d, %t, %s\n", v.Id, v.Credit, v.Debit, v.Store_id, v.Category_id, v.Applied, v.Trans_date.Format(time.RFC822))
  }
  budgetTotal, err := models.BudgetTotal(time.Now())
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("Balance: %d\n", budgetTotal)

  projBalance, err := models.ProjectedBalance()
  if err != nil {
    fmt.Println(err)
    return
  }
  enc := json.NewEncoder(os.Stdout)
  d := map[string]int{"projBalance": projBalance}
  enc.Encode(d)

  enc1 := json.NewEncoder(os.Stdout)
  budgetEntry, _ := models.BudgetEntry(1)
  enc1.Encode(budgetEntry)
  // fmt.Println(budgetEntry)
}
