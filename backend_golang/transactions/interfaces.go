package transactions

import "net/http"

type TransactionService interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
	AddTransaction(w http.ResponseWriter, r *http.Request)
	// UpdateTransaction(w http.ResponseWriter, r *http.Request)
	RemoveTransaction(w http.ResponseWriter, r *http.Request)
}
