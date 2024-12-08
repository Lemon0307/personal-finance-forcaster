package main

import (
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	database.InitDB()
	defer database.CloseDB()

	router := mux.NewRouter()
	budgets.BudgetRoutes(router, &budgets.BudgetHandler{})
	router.HandleFunc("/sign_up", auth.SignUpHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
