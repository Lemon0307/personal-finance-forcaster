package budgets

type Items struct {
	ItemName    string  `json:"item_name"`
	BudgetCost  float64 `json:"budget_cost"`
	Description string  `json:"description"`
	Priority    int32   `json:"priority"`
}

type Budget struct {
	BudgetName string   `json:"budget_name"`
	Items      []*Items `json:"items"`
}

type BudgetHandler struct{}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

type Response struct {
	Message    string
	StatusCode int
}
