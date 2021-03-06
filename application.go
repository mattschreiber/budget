package main

import (
	// "encoding/json"
	"github.com/mattschreiber/budget/models"
	// "fmt"
	// "time"
	// "os"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mattschreiber/budget/controllers"
	"github.com/mattschreiber/budget/middlewares"
	"github.com/mattschreiber/budget/tasks"
)

func main() {
	models.InitDB()

	go tasks.ScheduledTasks()

	r := mux.NewRouter()
	r.HandleFunc("/login", middlewares.LoginHandler).Methods("POST")
	r.HandleFunc("/projected-balance/{endDate}", middlewares.ValidateToken(controllers.GetProjBalance)).Methods("GET", "OPTIONS")
	r.HandleFunc("/budget-entries/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetBudgetEntries)).Methods("GET", "OPTIONS")
	r.HandleFunc("/ledger-entries/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetLedgerEntries)).Methods("GET", "OPTIONS")
	r.HandleFunc("/createledger", middlewares.ValidateToken(controllers.CreateLedgerEntry)).Methods("POST", "OPTIONS")
	r.HandleFunc("/createbudget", middlewares.ValidateToken(controllers.CreateBudgetEntry)).Methods("POST", "OPTIONS")
	r.HandleFunc("/balances/{startDate}/{endDate}", middlewares.ValidateToken(controllers.GetBalances)).Methods("GET", "OPTIONS")
	r.HandleFunc("/store-category", middlewares.ValidateToken(controllers.GetStoreCategory)).Methods("GET", "OPTIONS")
	r.HandleFunc("/createstore", middlewares.ValidateToken(controllers.CreateStore)).Methods("POST", "OPTIONS")
	r.HandleFunc("/ledger-entry/{id}", middlewares.ValidateToken(controllers.DeleteLedgerEntry)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/budget-entry/{id}", middlewares.ValidateToken(controllers.DeleteBudgetEntry)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/updateledger", middlewares.ValidateToken(controllers.UpdateLedgerEntry)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/updatebudget", middlewares.ValidateToken(controllers.UpdateBudgetEntry)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/reports/categories", middlewares.ValidateToken(controllers.GetAmountsByCategory)).Methods("GET", "OPTIONS").Queries("month", "{month}", "year", "{year}")
	r.HandleFunc("/reports/monthly/amounts", middlewares.ValidateToken(controllers.GetMonthlyTotalSpent)).Methods("GET", "OPTIONS").Queries("startDate", "{startDate}", "endDate", "{endDate}")
	r.HandleFunc("/payment-types", middlewares.ValidateToken(controllers.GetPaymentTypes)).Methods("GET", "OPTIONS")

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
