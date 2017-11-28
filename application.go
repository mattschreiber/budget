package main

import (
  "encoding/json"
  "budget/models"
  // "fmt"
  "time"
  "os"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "budget/middlewares"
  "budget/controllers"
)
func main() {
  models.InitDB()

  // budgetTotal, err := models.BudgetTotal(time.Now())
  // if err != nil {
  //   fmt.Println(err)
  //   return
  // }
  // fmt.Printf("Balance: %d\n", budgetTotal)

  enc1 := json.NewEncoder(os.Stdout)
  // budgetEntry, _ := models.BudgetEntry(1)
  // enc1.Encode(budgetEntry)
  // fmt.Println(budgetEntry)
  beginningOfTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
  endOfTime := time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC)

  amounts, _ := models.GetAmountSpent(beginningOfTime, endOfTime)
  enc1.Encode(amounts)
  // categories, _ := models.AllCategories()
  // for _, v := range categories {
  //   fmt.Printf("%d, %s\n", v.Id, v.Category_name)
  // }

  r := mux.NewRouter()
  r.HandleFunc("/login", middlewares.LoginHandler).Methods("POST")
  r.HandleFunc("/projected-balance/{endDate}", middlewares.ValidateToken(controllers.GetProjBalance)).Methods("GET", "OPTIONS")
  r.HandleFunc("/budget-entries/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetBudgetEntries)).Methods("GET", "OPTIONS")
  r.HandleFunc("/ledger-entries/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetLedgerEntries)).Methods("GET", "OPTIONS")
  r.HandleFunc("/balances/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetBalances)).Methods("GET", "OPTIONS")
  http.Handle("/", &MyServer{r})
  log.Fatal(http.ListenAndServe(":5000", nil))
}

type MyServer struct {
    r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    if origin := req.Header.Get("Origin"); origin != "" {
        rw.Header().Set("Access-Control-Allow-Origin", origin)
        rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        rw.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }
    // Stop here if its Preflighted OPTIONS request
    if req.Method == "OPTIONS" {
        return
    }
    // Lets Gorilla work
    s.r.ServeHTTP(rw, req)
}
