package main

import (
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"golang/forecast"
	"golang/routes"
	"golang/transactions"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	database.InitDB()
	defer database.CloseDB()

	router := mux.NewRouter()
	authRouter := router.PathPrefix("/auth").Subrouter()
	// separated routes from auth so that auth does not use
	// jwt middleware
	mainRouter := router.PathPrefix("/main").Subrouter()

	routes.AuthenticationRoutes(authRouter, &auth.AuthenticationHandler{})

	mainRouter.Use(auth.JWTAuthMiddleware)

	routes.BudgetRoutes(mainRouter, &budgets.BudgetHandler{})
	routes.TransactionRoutes(mainRouter, &transactions.TransactionHandler{})
	routes.ForecastRoutes(mainRouter, &forecast.ForecastHandler{})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Apply CORS to the router
	router_with_cors := c.Handler(router)

	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", router_with_cors)
	if err != nil {
		log.Fatal(err)
	}
}