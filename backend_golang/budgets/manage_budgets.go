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

// check if a budget with budget_name and user_id exists in the db
func BudgetExists(db *sql.DB, user_id string, budget_name string) bool {
	var res bool
	_ = db.QueryRow("SELECT EXISTS(SELECT * FROM Budget WHERE budget_name = ? AND user_id = ?)", budget_name, user_id).Scan(&res)
	return res
}

// check if an itme with item_name and user_id exists in the db
func ItemExists(db *sql.DB, user_id string, item_name string) bool {
	var res bool
	_ = db.QueryRow("SELECT EXISTS(SELECT * FROM Budget_Items WHERE item_name = ? AND user_id = ?)",
		item_name, user_id).Scan(&res)
	return res
}

func (budget *BudgetHandler) AddBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
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
	budget_exists := BudgetExists(database.DB, user_id, manageBudget.Budget.BudgetName)
	if !budget_exists {
		// add budget
		_, err := database.DB.Exec("INSERT INTO Budget (user_id, budget_name) VALUES (?, ?)",
			user_id, manageBudget.Budget.BudgetName)
		if err != nil {
			log.Fatal(err)
		}
		// add all budget items
		for i := 0; i < len(manageBudget.BudgetItems); i++ {
			// check if budget item exists
			var budget_item_exists bool
			err = database.DB.QueryRow(`SELECT EXISTS(SELECT * FROM Budget_Items WHERE item_name = ? 
			AND user_id = ? AND budget_name = ?)`, manageBudget.BudgetItems[i].ItemName, user_id,
				manageBudget.Budget.BudgetName).Scan(&budget_item_exists)

			if err != nil {
				log.Fatal(err)
			}

			if !budget_item_exists {
				_, err := database.DB.Exec(`INSERT INTO Budget_Items (user_id, budget_name,
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
			} else {
				response := ErrorMessage{
					Message: `Budget item with name ` + manageBudget.BudgetItems[i].ItemName +
						`already exists`,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
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
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID

	// query all budgets and its budget items (modified join to left join)
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
	// make a hash map to group budget with its budget items
	var budgets = make(map[string]*ManageBudgets)

	for rows.Next() {
		// add one row of budgets into structs
		var user_id string
		var budget Budget
		var budget_item BudgetItems
		if err := rows.Scan(&user_id, &budget.BudgetName, &budget_item.ItemName,
			&budget_item.BudgetCost, &budget_item.Description, &budget_item.Priority); err != nil {
			log.Fatal(err)
		}
		// add budget to the hash map if budget doesn't exist in the hash map
		if _, budget_exists := budgets[budget.BudgetName]; !budget_exists {
			budgets[budget.BudgetName] = &ManageBudgets{
				UserID: user_id,
				Budget: Budget{
					BudgetName: budget.BudgetName,
				},
				BudgetItems: []BudgetItems{},
			}
		}
		// group all budget items into one budget
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

	// group all budgets into an array of budgets
	var budgets_array []ManageBudgets
	for _, budget := range budgets {
		budgets_array = append(budgets_array, *budget)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets_array)
}

func (budget *BudgetHandler) RemoveBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID

	// check if budget exists
	if BudgetExists(database.DB, user_id, budget_name) {
		// delete budget from db
		_, err := database.DB.Exec(`DELETE FROM Budget WHERE user_id = ? AND budget_name = ?`, user_id, budget_name)
		if err != nil {
			log.Fatal(err)
		}
		// return success message
		response := Response{
			Message:    "Successfully deleted budget",
			StatusCode: 201,
		}
		// make json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// return error message
		response := ErrorMessage{
			Message:    "Cannot delete budget because budget does not exist",
			StatusCode: 401,
		}
		// make json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (budget *BudgetHandler) RemoveBudgetItems(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	item_name := vars["item_name"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID

	// check if item exists in the db
	if ItemExists(database.DB, user_id, item_name) {
		// delete item from db
		_, err := database.DB.Query(`DELETE FROM Budget_Items WHERE user_id = ? 
		AND budget_name = ? AND item_name = ?`, user_id, budget_name, item_name)
		if err != nil {
			log.Fatal(err)
		}
		// return success message
		response := Response{
			Message:    "Successfully deleted budget item",
			StatusCode: 201,
		}
		// make json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// return error message
		response := ErrorMessage{
			Message: "Cannot delete budget item because budget item does not exist",
		}
		// make json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func (budget *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID

	var manageBudget ManageBudgets
	err = json.NewDecoder(r.Body).Decode(&manageBudget)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	update_budget, err := database.DB.Exec(`UPDATE Budget SET budget_name = ? WHERE 
	budget_name = ? AND user_id = ?`, manageBudget.Budget.BudgetName, budget_name, user_id)
	fmt.Println("$1 Rows affected", update_budget)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{
		Message:    "Successfully updated budget",
		StatusCode: 201,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (budget *BudgetHandler) UpdateBudgetItems(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	item_name := vars["item_name"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID

	var budget_item BudgetItems
	_ = json.NewDecoder(r.Body).Decode(&budget_item)

	// building query string
	query := "UPDATE Budget_Items SET"
	args := []interface{}{}
	columns := []string{}

	// check if conditions exist
	if budget_item.BudgetCost != 0 {
		columns = append(columns, "budget_cost = ?")
		args = append(args, budget_item.BudgetCost)
	}
	if budget_item.Description != "" {
		columns = append(columns, "description = ?")
		args = append(args, budget_item.Description)
	}
	if budget_item.Priority != 0 {
		columns = append(columns, "priority = ?")
		args = append(args, budget_item.Priority)
	}

	if len(columns) > 0 {
		// complete query string
		query += " " + strings.Join(columns, ", ")
		query += " WHERE user_id = ? AND budget_name = ? AND item_name = ?"

		args = append(args, user_id, budget_name, item_name)

		_, err = database.DB.Exec(query, args...)
		if err != nil {
			log.Fatal(err)
		}

		// return success message
		response := Response{
			Message:    "Successfully updated budget item",
			StatusCode: 201,
		}
		// build json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// return error message
		response := ErrorMessage{
			Message:    "There are no updates to the budget item",
			StatusCode: 401,
		}

		// build json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}

// routing
func BudgetRoutes(router *mux.Router, budgetService BudgetService) {
	router.HandleFunc("/budgets", budgetService.GetBudget).Methods("GET")
	router.HandleFunc("/budgets/add_budget", budgetService.AddBudget).Methods("POST")
	router.HandleFunc("/budgets/update_budget/{budget_name}", budgetService.UpdateBudget).Methods("PUT")
	router.HandleFunc("/budgets/update_budget_item/{budget_name}/{item_name}", budgetService.UpdateBudgetItems).Methods("PUT")
	router.HandleFunc("/budgets/remove_budget/{budget_name}", budgetService.RemoveBudget).Methods("DELETE")
	router.HandleFunc("/budgets/remove_budget_item/{budget_name}/{item_name}", budgetService.RemoveBudgetItems).Methods("DELETE")
}
