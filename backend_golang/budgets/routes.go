package budgets

import (
	"github.com/gorilla/mux"
)

// routing
func BudgetRoutes(router *mux.Router, budgetService BudgetService) {
	router.HandleFunc("/budgets", budgetService.GetBudget).Methods("GET")
	router.HandleFunc("/budgets/add_budget", budgetService.AddBudget).Methods("POST")
	router.HandleFunc("/budgets/add_budget_item/{budget_name}", budgetService.AddItem).Methods("POST")
	router.HandleFunc("/budgets/update_budget/{budget_name}", budgetService.UpdateBudget).Methods("PUT")
	router.HandleFunc("/budgets/update_budget_item/{budget_name}/{item_name}", budgetService.UpdateItem).Methods("PUT")
	router.HandleFunc("/budgets/remove_budget/{budget_name}", budgetService.RemoveBudget).Methods("DELETE")
	router.HandleFunc("/budgets/remove_budget_item/{budget_name}/{item_name}", budgetService.RemoveItem).Methods("DELETE")
}
