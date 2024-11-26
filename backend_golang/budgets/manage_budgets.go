package budgets

import (
	"log"
	"net/http"
)

func AddBudget() {

}

func RemoveBudget() {

}

func UpdateBudget() {

}

func GetBudget() {

}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "budgets/":
		GetBudget()
	case "budgets/add_budget":
		AddBudget()
	case "budgets/remove_budget":
		RemoveBudget()
	case "budgets/update_budget":
		UpdateBudget()
	default:
		log.Fatal("Something went wrong with the URL.")
	}
}
