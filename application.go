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
  "budget/controllers"
)
func main() {
  models.InitDB()

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
  r.HandleFunc("/login", middlewares.LoginHandler).Methods("POST")
  r.HandleFunc("/home/{endDate}", middlewares.ValidateToken(controllers.GetProjBalance)).Methods("GET", "OPTIONS")
  r.HandleFunc("/balance/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetBalanceEntries)).Methods("GET", "OPTIONS")
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
