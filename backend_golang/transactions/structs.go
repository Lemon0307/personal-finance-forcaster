package transactions

import "time"

// type BudgetItem struct {
// 	ItemName    string  `json:"item_name"`
// 	BudgetCost  float32 `json:"budget_cost"`
// 	Description string  `json:"description"`
// 	Priority    float32 `json:"priority"`
// }

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
	ItemName     string         `json:"item_name"`
	Transactions []Transactions `json:"transactions"`
	MonthlyCosts MonthlyCosts   `json:"monthly_costs"`
}

type Date struct {
	time.Time
}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

type Response struct {
	Message    string
	StatusCode int
}
