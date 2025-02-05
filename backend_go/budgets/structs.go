package budgets

type Budget struct {
	BudgetName string `json:"budget_name"`
}

type Items struct {
	ItemName    string  `json:"item_name"`
	BudgetCost  float64 `json:"budget_cost"`
	Description string  `json:"description"`
	Priority    int32   `json:"priority"`
}

type B_Session struct {
	Budget Budget   `json:"budget"`
	Items  []*Items `json:"items"`
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
