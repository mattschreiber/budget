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
  startDate, _ := time.Parse(layout, params["startDate"])
  endDate, _ := time.Parse(layout, params["endDate"])

  categoryAmounts, err := models.AmountsByCategory(startDate, endDate)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Error in request")
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(categoryAmounts)
}
