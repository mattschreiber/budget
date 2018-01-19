package controllers

import (
  "encoding/json"
  "budget/models"
  "fmt"
  "time"
  "net/http"

  "github.com/gorilla/mux"
)

func GetAmountsByCategory(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  // startDate, _ := time.Parse(layout, params["startDate"])
  // endDate, _ := time.Parse(layout, params["endDate"])
  month := params["month"]
  year := params["year"]
  categoryAmounts, err := models.AmountsByCategory(month, year)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Error in request")
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(categoryAmounts)
}

func GetMonthlyTotalSpent(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  startDate, _ := time.Parse(layout, params["startDate"])
  endDate, _ := time.Parse(layout, params["endDate"])
  monthTotals, err := models.MonthlyTotalSpent(startDate, endDate)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Error in request")
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(monthTotals)
}
