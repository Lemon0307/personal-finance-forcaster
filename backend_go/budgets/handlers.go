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
	var budgets Budget
	// parse json into structs
	err = json.NewDecoder(r.Body).Decode(&budgets)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if budget exists in the database
	if BudgetExists(database.DB, user_id, budgets.BudgetName) {
		// return error message
		http.Error(w, "A budget with this name already exists.", http.StatusConflict)
	} else {
		// add budget
		_, err := database.DB.Exec("INSERT INTO Budget (user_id, budget_name) VALUES (?, ?)",
			user_id, budgets.BudgetName)
		if err != nil {
			log.Fatal(err)
		}
		// add all budget items
		for i := 0; i < len(budgets.Items); i++ {
			if !ItemExists(database.DB, user_id, budgets.Items[i].ItemName,
				budgets.BudgetName) { // add item if item does not exist
				_, err := database.DB.Exec(`INSERT INTO Budget_Items (user_id, budget_name,
				item_name, description, budget_cost, priority) VALUES (?, ?, ?, ?, ?, ?)`,
					user_id,
					budgets.BudgetName,
					budgets.Items[i].ItemName,
					budgets.Items[i].Description,
					budgets.Items[i].BudgetCost,
					budgets.Items[i].Priority)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				http.Error(w, `Budget item with name `+budgets.Items[i].ItemName+
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

	// create a new add budget item budget
	var item Items

	// parse json into structs
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if item already exists in the database
	if ItemExists(database.DB, user_id, item.ItemName, budget_name) {
		// return error message
		http.Error(w, "An item with this name already exists.", http.StatusConflict)
	} else {
		// add budget item to the database
		_, err = database.DB.Exec(`INSERT INTO Budget_Items (item_name, budget_name, user_id, 
		description, budget_cost, priority) VALUES (?, ?, ?, ?, ?, ?)`,
			item.ItemName,
			budget_name,
			user_id,
			item.Description,
			item.BudgetCost,
			item.Priority)
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
        SELECT budget.budget_name, 
               items.item_name, items.budget_cost, 
               items.description, items.priority
        FROM Budget budget
        LEFT JOIN Budget_Items items
        ON budget.user_id = items.user_id AND budget.budget_name = items.budget_name
        WHERE budget.user_id = ?`, user_id)
	if err != nil {
		log.Fatal(err)
	}

	var budgets = make(map[string]*Budget)

	for rows.Next() {
		var item Items
		// handle budget items and their possible null values
		var budget_name string
		var item_name sql.NullString
		var budget_cost sql.NullFloat64
		var description sql.NullString
		var priority sql.NullInt32

		// Scan values from the query result
		if err := rows.Scan(&budget_name, &item_name,
			&budget_cost, &description, &priority); err != nil {
			log.Fatal(err)
		}

		// Check if this budget has already been processed
		if _, budget_exists := budgets[budget_name]; !budget_exists {
			budgets[budget_name] = &Budget{
				BudgetName: budget_name,
				Items:      []*Items{},
			}
		}

		// check if any value in the budget item is null
		if !item_name.Valid || !budget_cost.Valid || !description.Valid || !priority.Valid {
			// append a null value into budget items
			budgets[budget_name].Items = append(budgets[budget_name].Items, nil)
		} else {
			// set variables to corresponding budget item values
			item.ItemName = item_name.String
			item.BudgetCost = budget_cost.Float64
			item.Description = description.String
			item.Priority = priority.Int32
			// append the object to the budget items array
			budgets[budget_name].Items = append(budgets[budget_name].Items, &Items{
				ItemName:    item.ItemName,
				BudgetCost:  item.BudgetCost,
				Description: item.Description,
				Priority:    item.Priority,
			})
		}
	}

	// Convert map to an array of Budget
	var budgets_array []Budget
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

	var budgets Budget
	err = json.NewDecoder(r.Body).Decode(&budgets)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// sql query to update budget in the database
	_, err = database.DB.Exec(`UPDATE Budget SET budget_name = ? WHERE 
	budget_name = ? AND user_id = ?`, budgets.BudgetName, budget_name, user_id)
	if err != nil {
		log.Fatal(err)
	}

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

	var items Items
	_ = json.NewDecoder(r.Body).Decode(&items)

	// building query string
	query := "UPDATE Budget_Items SET"
	args := []interface{}{}
	columns := []string{}

	// check if conditions exist
	if items.BudgetCost != 0 {
		columns = append(columns, "budget_cost = ?")
		args = append(args, items.BudgetCost)
	}
	if items.Description != "" {
		columns = append(columns, "description = ?")
		args = append(args, items.Description)
	}
	if items.Priority != 0 {
		columns = append(columns, "priority = ?")
		args = append(args, items.Priority)
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
