package transactions

import (
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/database"
	"log"
	"net/http"
	"strings"
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

func (transaction *TransactionHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	// vars := mux.Vars(r)
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID
	fmt.Println(user_id)

	var manageTransactions ManageTransactions
	json.NewDecoder(r.Body).Decode(&manageTransactions)

	month := MonthToInt(manageTransactions.Transactions[0].Date.Month().String())
	year := manageTransactions.Transactions[0].Date.Year()

	var exists bool
	err = database.DB.QueryRow(`SELECT EXISTS(SELECT * FROM Monthly_Costs WHERE 
	month = ? AND year = ? AND item_name = ? AND user_id = ?)`,
		month,
		year,
		manageTransactions.ItemName,
		user_id).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = database.DB.Exec(`INSERT INTO Monthly_Costs (user_id, item_name, month, year) VALUES (?, ?, ?, ?)`,
			user_id,
			manageTransactions.ItemName,
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

func (transaction *TransactionHandler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
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
}

func (transaction *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
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
}

func (transaction *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// var err error
	// vars := mux.Vars(r)
	// year := vars["year"]
	// month := vars["month"]
	// // check if token is valid or expired
	// token := r.Header.Get("Authorization")
	// token = strings.TrimPrefix(token, "Bearer ")
	// claims, err := auth.ValidateJWT(token)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // get user id from jwt
	// user_id := claims.UserID
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
	router.HandleFunc("/transactions/{year}/{month}/update_transaction/{transaction_id}", TransactionService.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{year}/{month}/remove_transaction/{transaction_id}", TransactionService.RemoveTransaction).Methods("DELETE")
}
