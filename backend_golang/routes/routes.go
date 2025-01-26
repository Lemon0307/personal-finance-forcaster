package routes

import (
	"golang/auth"
	"golang/budgets"
	"golang/forecast"
	"golang/transactions"

	"github.com/gorilla/mux"
)

func AuthenticationRoutes(router *mux.Router, authService auth.AuthenticationService) {
	router.HandleFunc("/sign_up", authService.SignUp).Methods("POST")
	router.HandleFunc("/login", authService.Login).Methods("POST")
}

func BudgetRoutes(router *mux.Router, budgetService budgets.BudgetService) {
	router.HandleFunc("/budgets", budgetService.GetBudget).Methods("GET")
	router.HandleFunc("/budgets/add_budget", budgetService.AddBudget).Methods("POST")
	router.HandleFunc("/budgets/add_budget_item/{budget_name}", budgetService.AddItem).Methods("POST")
	router.HandleFunc("/budgets/update_budget/{budget_name}", budgetService.UpdateBudget).Methods("PUT")
	router.HandleFunc("/budgets/update_budget_item/{budget_name}/{item_name}", budgetService.UpdateItem).Methods("PUT")
	router.HandleFunc("/budgets/remove_budget/{budget_name}", budgetService.RemoveBudget).Methods("DELETE")
	router.HandleFunc("/budgets/remove_budget_item/{budget_name}/{item_name}", budgetService.RemoveItem).Methods("DELETE")
}

func TransactionRoutes(router *mux.Router, TransactionService transactions.TransactionService) {
	router.HandleFunc("/transactions/{budget_name}/{item_name}/{year}/{month}",
		TransactionService.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{year}/{month}",
		TransactionService.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/add_transaction", TransactionService.AddTransaction).
		Methods("POST")
	router.HandleFunc(`/transactions/{year}/{month}/{budget_name}/{item_name}/remove_transaction/{transaction_id}`,
		TransactionService.RemoveTransaction).Methods("DELETE")
}

func ForecastRoutes(router *mux.Router, forecastService forecast.ForecastService) {
	router.HandleFunc("/forecast/forecast_transactions/{months}/{budget_name}/{item_name}",
		forecastService.ForecastTransactions).Methods("POST")
}
