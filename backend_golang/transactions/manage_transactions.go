package transactions

import (
	"log"
	"net/http"
)

func AddTransaction() {

}

func RemoveTransaction() {

}

func UpdateTransaction() {

}

func GetTransaction() {

}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "transactions/":
		GetTransaction()
	case "transactions/add_transaction":
		AddTransaction()
	case "transactions/remove_transaction":
		RemoveTransaction()
	case "transactions/update_transaction":
		UpdateTransaction()
	default:
		log.Fatal("Something went wrong with the URL.")
	}
}
