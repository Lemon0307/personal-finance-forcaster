package transactions

import (
	"golang/database"
	"golang/auth"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (transaction *TransactionHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID
}

func (transaction *TransactionHandler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	transaction_id := vars["transaction_id"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID
}

func (transaction *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	transaction_id := vars["transaction_id"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID
}

func (transaction *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	// get user id from jwt
	user_id := claims.UserID
}

func TransactionRoutes(router *mux.Router, TransactionService TransactionService) {
	router.HandleFunc("/transactions/{year}/{month}", TransactionService.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{year}/{month}/add_transaction", TransactionService.AddTransaction).Methods("POST")
	router.HandleFunc("/transactions/{year}/{month}/update_transaction/{transaction_id}", TransactionService.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{year}/{month}/remove_transaction/{transaction_id}", TransactionService.RemoveTransaction).Methods("DELETE")
}
