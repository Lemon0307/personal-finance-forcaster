package transactions

import "time"

type Item struct {
	BudgetName string `json:"budget_name"`
	ItemName   string `json:"item_name"`
}

type Transactions struct {
	TransactionID   string  `json:"transaction_id"`
	Date            Date    `json:"date"`
	TransactionType string  `json:"type"`
	TransactionName string  `json:"name"`
	Amount          float32 `json:"amount"`
}

type TransactionHandler struct{}

type MonthlyCosts struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type T_Session struct {
	Item         Item           `json:"item"`
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
