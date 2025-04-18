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
	router.HandleFunc("/get_current_balance",
		authService.GetCurrentBalance)
}

func BudgetRoutes(router *mux.Router, budgetService budgets.BudgetService) {
	router.HandleFunc("/budgets", budgetService.GetBudget).Methods("GET")
	router.HandleFunc("/budgets/add_budget", budgetService.AddBudget).Methods("POST")
	router.HandleFunc("/budgets/add_item/{budget_name}", budgetService.AddItem).Methods("POST")
	router.HandleFunc("/budgets/update_budget/{budget_name}", budgetService.UpdateBudget).Methods("PUT")
	router.HandleFunc("/budgets/update_item/{budget_name}/{item_name}", budgetService.UpdateItem).Methods("PUT")
	router.HandleFunc("/budgets/remove_budget/{budget_name}", budgetService.RemoveBudget).Methods("DELETE")
	router.HandleFunc("/budgets/remove_item/{budget_name}/{item_name}", budgetService.RemoveItem).Methods("DELETE")
}

func TransactionRoutes(router *mux.Router, transactionService transactions.TransactionService) {
	router.HandleFunc("/transactions/{budget_name}/{item_name}/{year}/{month}",
		transactionService.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{year}/{month}",
		transactionService.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/add_transaction", transactionService.AddTransaction).
		Methods("POST")
	router.HandleFunc(`/transactions/{year}/{month}/{budget_name}/{item_name}/remove_transaction/{transaction_id}`,
		transactionService.RemoveTransaction).Methods("DELETE")
}

func ForecastRoutes(router *mux.Router, forecastService forecast.ForecastService) {
	router.HandleFunc("/forecast/{months}/{budget_name}",
		forecastService.ForecastTransactions).Methods("GET")
}
