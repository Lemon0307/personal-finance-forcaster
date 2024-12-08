package budgets

import (
	"golang/auth"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (budget *BudgetHandler) AddBudget(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ValidateJWT(r.Header.Get("Authorization"))
	if err != nil {
		log.Fatal(err)
	}
	user_id := claims.UserID
	
	w.Header().Set("Content-Type", "application/json")
	response := ""
	w.Write([]byte(response))
}
func (budget *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {

}
func (budget *BudgetHandler) RemoveBudget(w http.ResponseWriter, r *http.Request) {

}
func (budget *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {

}

func BudgetRoutes(router *mux.Router, budgetService BudgetService) {
	router.HandleFunc("/budgets", budgetService.GetBudget).Methods("GET")
	router.HandleFunc("/budgets/remove_budget/{id}", budgetService.RemoveBudget).Methods("DELETE")
	router.HandleFunc("/budgets/add_budget", budgetService.AddBudget).Methods("POST")
	router.HandleFunc("/budgets/update_budget/{id}", budgetService.UpdateBudget).Methods("PUT")
}
