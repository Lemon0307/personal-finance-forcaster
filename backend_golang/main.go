package main

import (
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"golang/forecast"
	"golang/transactions"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	database.InitDB()
	defer database.CloseDB()

	router := mux.NewRouter()
	budgets.BudgetRoutes(router, &budgets.BudgetHandler{})
	transactions.TransactionRoutes(router, &transactions.TransactionHandler{})
	forecast.ForecastRoutes(router, &forecast.ForecastHandler{})
	router.HandleFunc("/sign_up", auth.SignUpHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/get_questions", auth.SQHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Apply CORS to the router
	router_with_cors := c.Handler(router)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", router_with_cors)
}
