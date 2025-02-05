package transactions

import (
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (transaction *TransactionHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// decodes json into session
	var session T_Session
	err = json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if item exists
	if budgets.ItemExists(database.DB, user_id, session.Item.ItemName,
		session.Item.BudgetName) {
		for i := 0; i < len(session.Transactions); i++ {
			// extracts month and year from the date of transaction
			month := MonthToInt(session.Transactions[i].Date.Month().String())
			year := session.Transactions[i].Date.Year()

			var exists bool
			// finds if a record of month and year exists in monthly costs table
			err = database.DB.QueryRow(`SELECT EXISTS(SELECT * FROM Monthly_Costs WHERE 
			month = ? AND year = ? AND item_name = ? AND user_id = ?)`,
				month,
				year,
				session.Item.ItemName,
				user_id).Scan(&exists)
			if err != nil {
				log.Fatal(err)
			}

			// if the record doesn't exist then add a record of monthly costs
			if !exists {
				_, err = database.DB.Exec(`INSERT INTO Monthly_Costs (user_id, item_name, 
				budget_name, month, year) VALUES (?, ?, ?, ?, ?)`,
					user_id,
					session.Item.ItemName,
					session.Item.BudgetName,
					month,
					year)
				if err != nil {
					log.Fatal(err)
				}
			}

			// generate a transaction id for the transaction
			session.Transactions[i].TransactionID = GenerateTransactionID()

			// add or subtract the current balance with transaction
			switch session.Transactions[i].TransactionType {
			case "inflow":
				_, err = database.DB.Exec(`UPDATE User SET current_balance = current_balance + ? 
				WHERE user_id = ?`, session.Transactions[0].Amount, user_id)
				if err != nil {
					log.Fatal(err)
				}
			case "outflow":
				_, err = database.DB.Exec(`UPDATE User SET current_balance = current_balance - ? 
				WHERE user_id = ?`, session.Transactions[0].Amount, user_id)
				if err != nil {
					log.Fatal(err)
				}
			}

			// add transaction data to the database
			_, err = database.DB.Exec(`INSERT INTO Transactions (user_id, transaction_id, 
			name, type, 
			amount, date, month, year, item_name, budget_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				user_id,
				session.Transactions[i].TransactionID,
				session.Transactions[i].TransactionName,
				session.Transactions[i].TransactionType,
				session.Transactions[i].Amount,
				session.Transactions[i].Date.Time,
				month,
				year,
				session.Item.ItemName,
				session.Item.BudgetName)
			if err != nil {
				log.Fatal(err)
			}
		}

		// success message
		response := Response{
			Message:    "Successfully added transaction",
			StatusCode: 201,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "The budget item does not exist, please try again",
			http.StatusNotFound)
	}
}

func (transaction *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	// get month, year, item name and budget name from url
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	item_name := vars["item_name"]
	budget_name := vars["budget_name"]

	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var results T_Session

	// find all transactions related to a budget item and a month, year
	rows, err := database.DB.Query(`
	SELECT 
		t.transaction_id,
		t.name,
		t.type,
		t.amount,
		t.date
	FROM 
		Transactions t
	WHERE 
		t.user_id = ?
		AND t.budget_name = ? 
		AND t.item_name = ? 
		AND t.month = ? 
		AND t.year = ?;
	`, user_id, budget_name, item_name, month, year)
	if err != nil {
		log.Fatal(err)
	}

	// add budget name, item name and month year to results
	results.Item.BudgetName = budget_name
	results.Item.ItemName = item_name
	results.MonthlyCosts = MonthlyCosts{Month: month, Year: year}

	// add array of transactions to response
	for rows.Next() {
		var t Transactions
		var date string

		err = rows.Scan(
			&t.TransactionID,
			&t.TransactionName,
			&t.TransactionType,
			&t.Amount,
			&date)
		if err != nil {
			log.Fatal(err)
		}

		// parsing date from db into format
		t.Date.Time, err = time.Parse("2006-01-02", date)
		if err != nil {
			log.Fatal(err)
		}

		// add data found into results
		results.Transactions = append(results.Transactions, t)
	}

	// show results
	w.Header().Set("Content-Type", "application/json")
	// parse results struct to json
	json.NewEncoder(w).Encode(results)
}

func (transaction *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	// get month, year from url
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])

	// Extract user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var results []T_Session

	// query all budget items of month and year
	rows, err := database.DB.Query(`
	SELECT 
		bi.budget_name,
		bi.item_name,
		mc.month,
		mc.year
	FROM 
		Budget_Items bi
	JOIN 
		Monthly_Costs mc
	ON
		bi.item_name = mc.item_name AND bi.budget_name = mc.budget_name AND bi.user_id = mc.user_id
	WHERE 
		bi.user_id = ? 
		AND mc.month = ? 
		AND mc.year = ?;
	`, user_id, month, year)
	if err != nil {
		log.Fatal(err)
	}

	// build results
	for rows.Next() {
		var bi Item
		var mc MonthlyCosts

		// get the budget name, item name, month and year from query row
		err = rows.Scan(&bi.BudgetName, &bi.ItemName, &mc.Month, &mc.Year)
		if err != nil {
			log.Fatal(err)
		}

		// get all transactions for the current budget of month and year
		var transactions []Transactions
		transRows, err := database.DB.Query(`
		SELECT 
			t.transaction_id,
			t.name,
			t.type,
			t.amount,
			t.date
		FROM 
			Transactions t
		WHERE 
			t.user_id = ? 
			AND t.month = ? 
			AND t.year = ? 
			AND t.budget_name = ? 
			AND t.item_name = ?;
		`, user_id, month, year, bi.BudgetName, bi.ItemName)
		if err != nil {
			log.Fatal(err)
		}

		// build transactions to append to budget
		for transRows.Next() {
			var t Transactions
			var date string

			err = transRows.Scan(&t.TransactionID, &t.TransactionName, &t.TransactionType, &t.Amount, &date)
			if err != nil {
				log.Fatal(err)
			}

			// Parsing date from db into format
			t.Date.Time, err = time.Parse("2006-01-02", date)
			if err != nil {
				log.Fatal(err)
			}

			transactions = append(transactions, t)
		}

		// Add the result for this budget item
		results = append(results, T_Session{
			Item:         bi,
			Transactions: transactions,
			MonthlyCosts: mc,
		})
	}

	// set transactions as nil for budgets with no transactions
	for i, item := range results {
		if len(item.Transactions) == 0 {
			results[i].Transactions = nil
		}
	}

	// show results
	w.Header().Set("Content-Type", "application/json")
	// parse results struct to json
	json.NewEncoder(w).Encode(results)
}

func (transaction *TransactionHandler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	// get month, year, item name and transaction id from url
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	item_name := vars["item_name"]
	budget_name := vars["budget_name"]
	transaction_id := vars["transaction_id"]

	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// check if transaction exists in the db
	transaction_exists := TransactionExists(database.DB, transaction_id, user_id, month, year,
		item_name, budget_name)

	// get amount and transaction type from transactions
	var t_type string
	var t_amount float32

	database.DB.QueryRow(`SELECT type, amount FROM Transactions 
	WHERE user_id = ? AND transaction_id = ?`, user_id, transaction_id).Scan(&t_type, &t_amount)

	// undos transaction addition/subtraction to current balance
	switch t_type {
	case "inflow":
		_, err = database.DB.Exec(`UPDATE User SET current_balance = current_balance - ? 
		WHERE user_id = ?`, t_amount, user_id)
		if err != nil {
			log.Fatal(err)
		}
	case "outflow":
		_, err = database.DB.Exec(`UPDATE User SET current_balance = current_balance + ? 
		WHERE user_id = ?`, t_amount, user_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	if transaction_exists {
		// delete transaction
		// checks if transaction belongs to the user
		_, err = database.DB.Exec(`
		DELETE FROM Transactions
		WHERE transaction_id = ? 
		AND user_id = (SELECT user_id FROM Monthly_Costs WHERE item_name = ? 
		AND month = ? AND year = ? AND budget_name = ?);
`,
			transaction_id, item_name, month, year, budget_name)
		if err != nil {
			log.Fatal(err)
		}
		// return success message
		response := Response{
			Message:    "Successfully deleted transaction",
			StatusCode: 201,
		}
		// return json response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// return error message
		http.Error(w, "transaction with id"+transaction_id+" and in budget item "+
			item_name+" doesn't exist, please try again", http.StatusNotFound)
	}
}
