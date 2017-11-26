package controllers

import (
  "encoding/json"
  "budget/models"
  "time"
  "net/http"

  "github.com/gorilla/mux"
)

func GetLedgerEntries(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  startDate, _ := time.Parse(layout, params["startDate"])
  endDate, _ := time.Parse(layout, params["endDate"])
  ledgerEntries, _ := models.AllLedgerEntries(startDate, endDate)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(ledgerEntries)
}
