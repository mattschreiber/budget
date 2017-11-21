package main

import (
  "encoding/json"
  "budget/models"
  "fmt"
  "time"
  "os"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "budget/middlewares"
)

func main() {
  models.InitDB()

  budgetEntries, _ := models.AllBudgetEntries()
  for _, v := range budgetEntries {
    fmt.Printf("%d, %d, %d, %d, %d, %t, %s, %s\n", v.Id, v.Credit, v.Debit, v.Store_id, v.Category_id, v.Applied, v.Trans_date.Format(time.RFC822), v.Store_name)
  }
  budgetTotal, err := models.BudgetTotal(time.Now())
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("Balance: %d\n", budgetTotal)

  enc1 := json.NewEncoder(os.Stdout)
  budgetEntry, _ := models.BudgetEntry(1)
  enc1.Encode(budgetEntry)
  // fmt.Println(budgetEntry)

  categories, _ := models.AllCategories()
  for _, v := range categories {
    fmt.Printf("%d, %s\n", v.Id, v.Category_name)
  }

  r := mux.NewRouter()
  r.HandleFunc("/home/{endDate}", middlewares.ValidateToken(GetProjBalance)).Methods("GET", "OPTIONS")
  r.HandleFunc("/login", middlewares.LoginHandler).Methods("POST")
  // r.HandleFunc("/valid", middlewares.ValidateToken).Methods("GET")
  http.Handle("/", &MyServer{r})
  log.Fatal(http.ListenAndServe(":5000", nil))

}

func GetProjBalance(w http.ResponseWriter, req *http.Request) {

  params := mux.Vars(req)
  fmt.Println(params["endDate"])
  layout := "2006-1-2"
  // startDate, _ := time.Parse(layout, "1900-01-01")
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
