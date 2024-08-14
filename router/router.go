package router

import (
	"eth-transaction-api/config"
	"eth-transaction-api/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

// NewRouter initializes and returns a new router
func NewRouter(cfg *config.Config, db *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	limiter := rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.BurstLimit)

	// Define the routes

	router.HandleFunc("/accounts/{accountUuid}/transactions", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		controllers.GetAccountTransactions(db)(w, r) // Pass the db instance to the controller
	}).Methods("GET")

	router.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		controllers.GetAccounts(db)(w, r) // Pass the db instance to the controller

	}).Methods("GET")

	router.HandleFunc("/pooled-eth-and-shares", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		controllers.GetPooledETHAndShares(db)(w, r) // Pass the db instance to the controller
	}).Methods("GET")

	router.HandleFunc("/last-5-depositors", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		controllers.GetLast5Depositors(db)(w, r) // Pass the db instance to the controller
	}).Methods("GET")

	return router
}
