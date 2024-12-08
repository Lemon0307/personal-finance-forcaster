package budgets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/database"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (budget *Budget) ValidateBudget(db *sql.DB, user_id string) (bool, error) {
	var res bool
	err := db.QueryRow("SELECT EXISTS(SELECT * FROM Budget WHERE budget_name = ? AND user_id = ?)", budget.BudgetName, user_id).Scan(&res)
	if err != nil {
		return false, err
	}
	return !res, nil
}

func (budget *BudgetHandler) AddBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID

	// create new manage budget session
	var manageBudget ManageBudgets
	// parse json into structs
	err = json.NewDecoder(r.Body).Decode(&manageBudget)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if budget exists in the database
	budget_ok, err := manageBudget.Budget.ValidateBudget(database.DB, user_id)
	if err != nil {
		log.Fatal(err)
	}
	if budget_ok {
		// add budget
		add_budget_query, err := database.DB.Exec("INSERT INTO Budget (user_id, budget_name) VALUES (?, ?)",
			user_id, manageBudget.Budget.BudgetName)
		if err != nil {
			log.Fatal(err)
		}
		// add all budget items
		for i := 0; i < len(manageBudget.BudgetItems); i++ {
			add_budget_items_query, err := database.DB.Exec(`INSERT INTO Budget_Items (user_id, budget_name,
			 item_name, description, budget_cost, priority) VALUES (?, ?, ?, ?, ?, ?)`,
				user_id,
				manageBudget.Budget.BudgetName,
				manageBudget.BudgetItems[i].ItemName,
				manageBudget.BudgetItems[i].Description,
				manageBudget.BudgetItems[i].BudgetCost,
				manageBudget.BudgetItems[i].Priority)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("$1 rows affected in Users Table, $2 rows affected in Security Questions table", add_budget_query, add_budget_items_query)
		}
		// return success message
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message:    "Successfully added budget",
			StatusCode: 201,
		}
		// builds json response
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
			return
		}
	} else {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		response := ErrorMessage{
			Message:    "A budget with this name already exists.",
			StatusCode: 400,
		}
		// builds json response
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
			return
		}
	}
}
func (budget *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {

	var err error
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID

	// create new manage budget session (modified join to left join)
	rows, err := database.DB.Query(`
        SELECT budget.user_id, budget.budget_name, 
               budget_items.item_name, budget_items.budget_cost, 
			   budget_items.description, budget_items.priority
        FROM Budget budget
        LEFT JOIN Budget_Items budget_items
        ON budget.user_id = budget_items.user_id AND budget.budget_name = budget_items.budget_name
        WHERE budget.user_id = ?`, user_id)
	if err != nil {
		log.Fatal(err)
	}
	var budgets = make(map[string]*ManageBudgets)

	for rows.Next() {
		var user_id string
		var budget Budget
		var budget_item BudgetItems
		if err := rows.Scan(&user_id, &budget.BudgetName, &budget_item.ItemName,
			&budget_item.BudgetCost, &budget_item.Description, &budget_item.Priority); err != nil {
			log.Fatal(err)
		}

		if _, exists := budgets[budget.BudgetName]; !exists {
			budgets[budget.BudgetName] = &ManageBudgets{
				UserID: user_id,
				Budget: Budget{
					BudgetName: budget.BudgetName,
				},
				BudgetItems: []BudgetItems{},
			}
		}

		if budget_item.ItemName != "" {
			budgets[budget.BudgetName].BudgetItems = append(budgets[budget.BudgetName].BudgetItems, BudgetItems{
				ItemName:    budget_item.ItemName,
				BudgetCost:  budget_item.BudgetCost,
				Description: budget_item.Description,
				Priority:    budget_item.Priority,
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	var budgets_array []ManageBudgets
	for _, budget := range budgets {
		budgets_array = append(budgets_array, *budget)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets_array)
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
