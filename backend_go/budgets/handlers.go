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

func (budget *BudgetHandler) AddBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// create new manage budget session
	var manageBudget B_Session
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
		for i := 0; i < len(manageBudget.Items); i++ {
			// check if budget item exists

			if !ItemExists(database.DB, user_id, manageBudget.Items[i].ItemName,
				manageBudget.Budget.BudgetName) {
				_, err := database.DB.Exec(`INSERT INTO Budget_Items (user_id, budget_name,
				item_name, description, budget_cost, priority) VALUES (?, ?, ?, ?, ?, ?)`,
					user_id,
					manageBudget.Budget.BudgetName,
					manageBudget.Items[i].ItemName,
					manageBudget.Items[i].Description,
					manageBudget.Items[i].BudgetCost,
					manageBudget.Items[i].Priority)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				http.Error(w, `Budget item with name `+manageBudget.Items[i].ItemName+
					`already exists`, http.StatusConflict)
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
		http.Error(w, "A budget with this name already exists.", http.StatusConflict)
	}
}

func (budget *BudgetHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}
	budget_name := vars["budget_name"]

	// create a new add budget item session
	var session Items

	// parse json into structs
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if item already exists in the database
	if !ItemExists(database.DB, user_id, session.ItemName, budget_name) {
		// add budget item to the database
		_, err = database.DB.Exec(`INSERT INTO Budget_Items (item_name, budget_name, user_id, 
		description, budget_cost, priority) VALUES (?, ?, ?, ?, ?, ?)`,
			session.ItemName,
			budget_name,
			user_id,
			session.Description,
			session.BudgetCost,
			session.Priority)
		if err != nil {
			log.Fatal(err)
		}
		// return success message
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message:    "Successfully added budget item",
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
		http.Error(w, "A budget item with this name already exists.", http.StatusConflict)
	}
}

func (budget *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// Query to get budgets with associated items
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

	var budgets = make(map[string]*B_Session)

	for rows.Next() {
		var user_id string
		var budget Budget
		var budget_item Items
		// handle budget items and their possible null values
		var itemName sql.NullString
		var budgetCost sql.NullFloat64
		var description sql.NullString
		var priority sql.NullInt32

		// Scan values from the query result
		if err := rows.Scan(&user_id, &budget.BudgetName, &itemName,
			&budgetCost, &description, &priority); err != nil {
			log.Fatal(err)
		}

		// Check if this budget has already been processed
		if _, budget_exists := budgets[budget.BudgetName]; !budget_exists {
			budgets[budget.BudgetName] = &B_Session{
				Budget: Budget{
					BudgetName: budget.BudgetName,
				},
				Items: []*Items{},
			}
		}

		// check if any value in the budget item is null
		if !itemName.Valid || !budgetCost.Valid || !description.Valid || !priority.Valid {
			// append a null value into budget items
			budgets[budget.BudgetName].Items = append(budgets[budget.BudgetName].Items, nil)
		} else {
			// set variables to corresponding budget item values
			budget_item.ItemName = itemName.String
			budget_item.BudgetCost = budgetCost.Float64
			budget_item.Description = description.String
			budget_item.Priority = priority.Int32
			// append the object to the budget items array
			budgets[budget.BudgetName].Items = append(budgets[budget.BudgetName].Items, &Items{
				ItemName:    budget_item.ItemName,
				BudgetCost:  budget_item.BudgetCost,
				Description: budget_item.Description,
				Priority:    budget_item.Priority,
			})
		}
	}

	// Convert map to an array of B_Session
	var budgets_array []B_Session
	for _, budget := range budgets {
		budgets_array = append(budgets_array, *budget)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets_array)
}

func (budget *BudgetHandler) RemoveBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// check if budget exists
	if BudgetExists(database.DB, user_id, budget_name) {
		// delete budget from db
		_, err = database.DB.Exec(`DELETE FROM Budget WHERE user_id = ? AND budget_name = ?`, user_id, budget_name)
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
		http.Error(w, "Cannot delete budget because budget does not exist", http.StatusNotFound)
	}
}

func (budget *BudgetHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	item_name := vars["item_name"]
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// check if item exists in the db
	if ItemExists(database.DB, user_id, item_name, budget_name) {
		// delete item from db
		_, err = database.DB.Query(`DELETE FROM Budget_Items WHERE user_id = ? 
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
		http.Error(w, "Cannot delete budget item because budget item does not exist", http.StatusNotFound)
	}
}
func (budget *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var manageBudget B_Session
	err = json.NewDecoder(r.Body).Decode(&manageBudget)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// sql query to update budget in the database
	update_budget, err := database.DB.Exec(`UPDATE Budget SET budget_name = ? WHERE 
	budget_name = ? AND user_id = ?`, manageBudget.Budget.BudgetName, budget_name, user_id)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := update_budget.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)

	// success message
	response := Response{
		Message:    "Successfully updated budget",
		StatusCode: 201,
	}
	// encode to a json response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (budget *BudgetHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	budget_name := vars["budget_name"]
	item_name := vars["item_name"]
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var budget_item Items
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
		http.Error(w, "There are no updates to the budget item", http.StatusNotModified)
	}

}
