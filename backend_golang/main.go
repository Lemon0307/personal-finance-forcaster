package main

import (
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
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
	router.HandleFunc("/sign_up", auth.SignUpHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")

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
