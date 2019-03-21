package controllers

import (
  "encoding/json"
  "github.com/mattschreiber/budget/models"
  "time"
  "net/http"
  "fmt"

  "github.com/gorilla/mux"
)

func GetLedgerEntries(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  startDate, _ := time.Parse(layout, params["startDate"])
  endDate, _ := time.Parse(layout, params["endDate"])
  ledgerEntries, err := models.AllLedgerEntries(startDate, endDate)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error in request")
		return
  }
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

func DeleteLedgerEntry(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  id := params["id"]
  count, err := models.DeleteLedgerEntry(id)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Error in request")
    return
  }
  deletedId := map[string]int64 {"count": count }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(deletedId)
}
