package controllers

import (
  "encoding/json"
  "budget/models"
  "fmt"
  "time"
  "net/http"

  "github.com/gorilla/mux"
)

const layout = "2006-1-2"

type TokenContextKey struct {
  Name string
}

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

func GetStoreCategory(w http.ResponseWriter, req *http.Request) {
  storeCategory, _ := models.GetStoreCategory()
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(storeCategory)
}

func CreateBudgetEntry(w http.ResponseWriter, req *http.Request){
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

func DeleteBudgetEntry(w http.ResponseWriter, req *http.Request) {
  params := mux.Vars(req)
  id := params["id"]
  count, err := models.DeleteBudgetEntry(id)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, "Error in request")
    return
  }
  deletedId := map[string]int64 {"count": count }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(deletedId)
}

func CreateStore(w http.ResponseWriter, req *http.Request){
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
