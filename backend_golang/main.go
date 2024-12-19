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
	// separated routes from auth so that auth does not use
	// jwt middleware
	mainRouter := router.PathPrefix("/main").Subrouter()

	auth.AuthenticationRoutes(router, &auth.AuthenticationHandler{})

	mainRouter.Use(auth.JWTAuthMiddleware)

	budgets.BudgetRoutes(mainRouter, &budgets.BudgetHandler{})
	transactions.TransactionRoutes(mainRouter, &transactions.TransactionHandler{})
	forecast.ForecastRoutes(mainRouter, &forecast.ForecastHandler{})
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
