package forecast

type ForecastHandler struct{}

type TotalTransactions struct {
	Amount string
	Month  string
	Year   string
}
type Item struct {
	ItemName          string
	BudgetCost        float64
	TotalTransactions []TotalTransactions
}

type Items struct {
	BudgetName string
	Items      []Item
}
