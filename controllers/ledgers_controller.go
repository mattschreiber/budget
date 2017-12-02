package controllers

import (
  "encoding/json"
  "budget/models"
  "time"
  "net/http"
  "fmt"

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

func CreateLedgerEntry(w http.ResponseWriter, req *http.Request){
  // params := mux.Vars(req)
  var ledger models.Model
  id := make(map[string]int)
  err := json.NewDecoder(req.Body).Decode(&ledger)
  if err != nil {
    fmt.Println(err)
    http.Error(w, http.StatusText(500), 500)
    return
  }
  id["id"], err = models.CreateLedgerEntry(ledger)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(id)
}
