package budgets

import "net/http"

type BudgetService interface {
	AddBudget(w http.ResponseWriter, r *http.Request)
	GetBudget(w http.ResponseWriter, r *http.Request)
	RemoveBudget(w http.ResponseWriter, r *http.Request)
	UpdateBudget(w http.ResponseWriter, r *http.Request)
}