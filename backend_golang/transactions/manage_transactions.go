package transactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GenerateTransactionID() string {
	res := uuid.New()
	return res.String()
}

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

func TransactionExists(db *sql.DB, transaction_id string, user_id string) bool {
	var res bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM transactions WHERE transaction_id = ?
	AND user_id = ?)`, transaction_id, user_id).
		Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
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

	var manageTransactions ManageTransactions
	err = json.NewDecoder(r.Body).Decode(&manageTransactions)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	month := MonthToInt(manageTransactions.Transactions[0].Date.Month().String())
	year := manageTransactions.Transactions[0].Date.Year()

	var exists bool
	err = database.DB.QueryRow(`SELECT EXISTS(SELECT * FROM Monthly_Costs WHERE 
	month = ? AND year = ? AND item_name = ? AND user_id = ?)`,
		month,
		year,
		manageTransactions.BudgetItem.ItemName,
		user_id).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = database.DB.Exec(`INSERT INTO Monthly_Costs (user_id, item_name, month, year) VALUES (?, ?, ?, ?)`,
			user_id,
			manageTransactions.BudgetItem.ItemName,
			month,
			year)
		if err != nil {
			log.Fatal(err)
		}
	}

	manageTransactions.Transactions[0].TransactionID = GenerateTransactionID()
	_, err = database.DB.Exec(`INSERT INTO Transactions (user_id, transaction_id, transaction_name, transaction_type, 
	amount, date, month, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user_id,
		manageTransactions.Transactions[0].TransactionID,
		manageTransactions.Transactions[0].TransactionName,
		manageTransactions.Transactions[0].TransactionType,
		manageTransactions.Transactions[0].Amount,
		manageTransactions.Transactions[0].Date.Time,
		month,
		year)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{
		Message:    "Successfully added transaction",
		StatusCode: 201,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func (transaction *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
// var err error
// vars := mux.Vars(r)
// year := vars["year"]
// month := vars["month"]
// transaction_id := vars["transaction_id"]
// // check if token is valid or expired
// token := r.Header.Get("Authorization")
// token = strings.TrimPrefix(token, "Bearer ")
// claims, err := auth.ValidateJWT(token)
// if err != nil {
// 	log.Fatal(err)
// }
// // get user id from jwt
// user_id := claims.UserID
// }

func (transaction *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])

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
			transaction.transaction_id,
			transaction.transaction_name,
			transaction.transaction_type,
			transaction.amount,
			transaction.date,
			budget_item.item_name,
			budget.budget_name
		FROM 
			Transactions transaction
		JOIN 
			Monthly_Costs monthly_costs ON transaction.user_id = monthly_costs.user_id AND transaction.month = 
			monthly_costs.month AND transaction.year = monthly_costs.year
		JOIN 
			Budget_Items budget_item ON monthly_costs.item_name = budget_item.item_name
		JOIN 
			Budget budget ON budget_item.budget_name = budget.budget_name
		WHERE 
			transaction.month = ? AND transaction.year = ? AND transaction.user_id = ?;
	`, month, year, user_id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var data ManageTransactions
		var t Transactions
		var bi BudgetItem
		var date string

		err = rows.Scan(
			&t.TransactionID,
			&t.TransactionName,
			&t.TransactionType,
			&t.Amount,
			&date,
			&bi.ItemName,
			&bi.BudgetName)
		if err != nil {
			log.Fatal(err)
		}

		// parsing date from db into format
		t.Date.Time, err = time.Parse("2006-01-02", date)
		if err != nil {
			log.Fatal(err)
		}

		// add data found into results
		data.BudgetItem = bi
		data.Transactions = append(data.Transactions, t)
		data.MonthlyCosts = MonthlyCosts{Month: month, Year: year}

		results = data
	}

	// show results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (transaction *TransactionHandler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	item_name := vars["item_name"]
	transaction_id := vars["transaction_id"]

	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	// check if transaction exists in the db
	transaction_exists := TransactionExists(database.DB, transaction_id, user_id)
	if transaction_exists {
		// delete transaction
		// checks if transaction belongs to the user
		_, err = database.DB.Exec(`
		DELETE FROM Transactions
		WHERE transaction_id = ? 
		AND user_id = (SELECT user_id FROM Monthly_Costs WHERE item_name = ? AND month = ? AND year = ?);
`,
			transaction_id, item_name, month, year)
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
		http.Error(w, "transaction with id"+transaction_id+" and in budget item"+
			"item_name"+"doesn't exist, please try again", http.StatusNotFound)
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
	router.HandleFunc("/transactions/{year}/{month}", TransactionService.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/add_transaction", TransactionService.AddTransaction).Methods("POST")
	// router.HandleFunc("/transactions/{year}/{month}/update_transaction/{transaction_id}", TransactionService.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{year}/{month}/{item_name}/remove_transaction/{transaction_id}", TransactionService.RemoveTransaction).Methods("DELETE")
}
