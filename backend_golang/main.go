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

	router := mux.NewRouter()
	budgets.BudgetRoutes(router, &budgets.BudgetHandler{})
	http.HandleFunc("/sign_up", auth.SignUpHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
