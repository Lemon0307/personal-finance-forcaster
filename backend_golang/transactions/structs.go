package transactions

import "time"

type BudgetFeatures struct {
	BudgetID   string
	FeatureID  string
	CostBudget float32
	Priority   float32
}

type Transactions struct {
	TransactionID   string
	TransactionName string
	TransactionType string
	Amount          float32
	Date            Date
	Month           int
	Year            int
}

type MonthlyCosts struct {
	BudgetID string
	Month    int
	Year     int
}

type Date struct {
	time.Time
}
