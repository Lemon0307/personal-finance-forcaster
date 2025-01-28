package transactions

import "net/http"

type TransactionService interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
	GetAllTransactions(w http.ResponseWriter, r *http.Request)
	AddTransaction(w http.ResponseWriter, r *http.Request)
	RemoveTransaction(w http.ResponseWriter, r *http.Request)
	GetCurrentBalance(w http.ResponseWriter, r *http.Request)
}
