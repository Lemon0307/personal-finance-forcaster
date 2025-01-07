package transactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// generates transaction id
func GenerateTransactionID() string {
	res := uuid.New()
	return res.String()
}

// switch case to map month in string to map in int
func MonthToInt(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	}
	return 0
}

func TransactionExists(db *sql.DB, transaction_id, user_id string, month, year int,
	item_name, budget_name string) bool {
	var res bool
	// query transactions table for a record that contains user id and transaction id
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM transactions WHERE transaction_id = ?
	AND user_id = ? AND month = ? AND year = ? AND item_name = ? AND budget_name = ?)`,
		transaction_id, user_id, month, year, item_name, budget_name).
		Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return res
}

func (transaction *TransactionHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// decodes json into session
	var manageTransactions ManageTransactions
	err = json.NewDecoder(r.Body).Decode(&manageTransactions)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// check if budget item exists
	if budgets.ItemExists(database.DB, user_id, manageTransactions.BudgetItem.ItemName,
		manageTransactions.BudgetItem.BudgetName) {
		// extracts month and year from the date of transaction
		month := MonthToInt(manageTransactions.Transactions[0].Date.Month().String())
		year := manageTransactions.Transactions[0].Date.Year()

		var exists bool
		// finds if a record of month and year exists in monthly costs table
		err = database.DB.QueryRow(`SELECT EXISTS(SELECT * FROM Monthly_Costs WHERE 
		month = ? AND year = ? AND item_name = ? AND user_id = ?)`,
			month,
			year,
			manageTransactions.BudgetItem.ItemName,
			user_id).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}

		// if the record doesn't exist then add a record of monthly costs
		if !exists {
			_, err = database.DB.Exec(`INSERT INTO Monthly_Costs (user_id, item_name, 
			budget_name, month, year) VALUES (?, ?, ?, ?, ?)`,
				user_id,
				manageTransactions.BudgetItem.ItemName,
				manageTransactions.BudgetItem.BudgetName,
				month,
				year)
			if err != nil {
				log.Fatal(err)
			}
		}

		// generate a transaction id for the transaction
		manageTransactions.Transactions[0].TransactionID = GenerateTransactionID()
		// add transaction data to the database
		_, err = database.DB.Exec(`INSERT INTO Transactions (user_id, transaction_id, 
		transaction_name, transaction_type, 
		amount, date, month, year, item_name, budget_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user_id,
			manageTransactions.Transactions[0].TransactionID,
			manageTransactions.Transactions[0].TransactionName,
			manageTransactions.Transactions[0].TransactionType,
			manageTransactions.Transactions[0].Amount,
			manageTransactions.Transactions[0].Date.Time,
			month,
			year,
			manageTransactions.BudgetItem.ItemName,
			manageTransactions.BudgetItem.BudgetName)
		if err != nil {
			log.Fatal(err)
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

	var results ManageTransactions

	// find all transactions related to a budget item and a month, year
	rows, err := database.DB.Query(`
	SELECT 
		t.transaction_id,
		t.transaction_name,
		t.transaction_type,
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
	results.BudgetItem.BudgetName = budget_name
	results.BudgetItem.ItemName = item_name
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

	fmt.Println(results)
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

	var results []ManageTransactions

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
		var bi BudgetItem
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
			t.transaction_name,
			t.transaction_type,
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
		results = append(results, ManageTransactions{
			BudgetItem:   bi,
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

	fmt.Println(results)

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
		fmt.Println(transaction_exists)
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

// parse date in the right format
func (date *Date) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])

	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	date.Time = parsedTime
	return nil
}

func TransactionRoutes(router *mux.Router, TransactionService TransactionService) {
	router.HandleFunc("/transactions/{budget_name}/{item_name}/{year}/{month}",
		TransactionService.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{year}/{month}",
		TransactionService.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/add_transaction", TransactionService.AddTransaction).
		Methods("POST")
	router.HandleFunc(`/transactions/{year}/{month}/{budget_name}/{item_name}/remove_transaction/{transaction_id}`,
		TransactionService.RemoveTransaction).Methods("DELETE")
}
