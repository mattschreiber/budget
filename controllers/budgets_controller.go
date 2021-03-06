package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mattschreiber/budget/models"

	"github.com/gorilla/mux"
)

const layout = "2006-1-2"

// TokenContextKey ...
type TokenContextKey struct {
	Name string
}

// GetProjBalance ...
func GetProjBalance(w http.ResponseWriter, req *http.Request) {

	t := TokenContextKey{}
	if token := req.Context().Value(t.Name); token != nil {
		// User is logged in
	} else {
		// User is not logged in
	}
	params := mux.Vars(req)
	endDate, _ := time.Parse(layout, params["endDate"])

	projBalance, err := models.ProjectedBalance(endDate)
	if err != nil {
		fmt.Println(err)
		return
	}
	// enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"projBalance": projBalance}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)

}

// GetBudgetEntries ...
func GetBudgetEntries(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	startDate, _ := time.Parse(layout, params["startDate"])
	endDate, _ := time.Parse(layout, params["endDate"])
	budgetEntries, err := models.AllBudgetEntries(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error in request")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgetEntries)
}

// GetBalances ...
func GetBalances(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	startDate, _ := time.Parse(layout, params["startDate"])
	endDate, _ := time.Parse(layout, params["endDate"])
	balances, err := models.GetAmountSpent(startDate, endDate)
	if err != nil {
		fmt.Println("Error get balances: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
}

// GetStoreCategory ...
func GetStoreCategory(w http.ResponseWriter, req *http.Request) {
	storeCategory, _ := models.GetStoreCategory()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storeCategory)
}

// CreateBudgetEntry add new budget entry
func CreateBudgetEntry(w http.ResponseWriter, req *http.Request) {
	// params := mux.Vars(req)
	var budget models.Model
	id := make(map[string]int)
	err := json.NewDecoder(req.Body).Decode(&budget)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	id["id"], err = models.CreateBudgetEntry(budget)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

// DeleteBudgetEntry delete a budget entry
func DeleteBudgetEntry(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	count, err := models.DeleteBudgetEntry(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error in request")
		return
	}
	deletedId := map[string]int64{"count": count}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedId)
}

// UpdateBudgetEntry method to update a single budget entry
func UpdateBudgetEntry(w http.ResponseWriter, req *http.Request) {
	var budgetEntry models.Model
	err := json.NewDecoder(req.Body).Decode(&budgetEntry)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	count, err := models.UpdateBudgetEntry(budgetEntry)
	if err != nil {
		fmt.Println("Error parsing ledger entry: ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)

}

// CreateStore create a new store
func CreateStore(w http.ResponseWriter, req *http.Request) {
	var store models.Store
	id := make(map[string]int)
	err := json.NewDecoder(req.Body).Decode(&store)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	id["id"], err = models.CreateStore(store)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

// GetPaymentTypes used to return all payment types
func GetPaymentTypes(w http.ResponseWriter, req *http.Request) {
	paymentTypes, _ := models.GetPaymentTypes()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentTypes)
}
