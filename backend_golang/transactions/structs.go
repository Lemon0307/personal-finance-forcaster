package transactions

import "time"

type BudgetItem struct {
	BudgetName string `json:"budget_name"`
	ItemName   string `json:"item_name"`
}

type Transactions struct {
	TransactionID   string  `json:"transaction_id"`
	TransactionName string  `json:"transaction_name"`
	TransactionType string  `json:"transaction_type"`
	Amount          float32 `json:"amount"`
	Date            Date    `json:"date"`
}

type TransactionHandler struct{}

type MonthlyCosts struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type ManageTransactions struct {
	BudgetItem   BudgetItem     `json:"budget_item"`
	Transactions []Transactions `json:"transactions"`
	MonthlyCosts MonthlyCosts   `json:"monthly_costs"`
}

type Date struct {
	time.Time
}

type Response struct {
	Message    string
	StatusCode int
}
