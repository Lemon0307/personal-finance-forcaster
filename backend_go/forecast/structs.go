package forecast

type ForecastHandler struct{}

type TotalTransactions struct {
	Amount float64
	Month  int
	Year   int
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
