package forecast

import "time"

type ForecastHandler struct{}

type TotalTransactions struct {
	Amount float64
	Date   string
}
type Item struct {
	ItemName    string
	BudgetCost  float64
	TotalSpent  []TotalTransactions
	TotalEarned []TotalTransactions
}

type Items struct {
	BudgetName string
	Items      []Item
}

type Date struct {
	time.Time
}
