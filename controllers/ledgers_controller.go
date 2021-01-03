package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mattschreiber/budget/models"

	"github.com/gorilla/mux"
)

// GetLedgerEntries function that fetches all ledger entries between two dates
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

// CreateLedgerEntry function that creates a new ledger entry
func CreateLedgerEntry(w http.ResponseWriter, req *http.Request) {
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

// DeleteLedgerEntry function that deletes a single ledger entry
func DeleteLedgerEntry(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	count, err := models.DeleteLedgerEntry(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error in request")
		return
	}
	deletedId := map[string]int64{"count": count}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)
}

// UpdateLedgerEntry method to update a single ledger entry
func UpdateLedgerEntry(w http.ResponseWriter, req *http.Request) {
	var ledgerEntry models.Model
	err := json.NewDecoder(req.Body).Decode(&ledgerEntry)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	count, err := models.UpdateLedgerEntry(ledgerEntry)
	if err != nil {
		fmt.Println("Error parsing ledger entry: ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)

}
