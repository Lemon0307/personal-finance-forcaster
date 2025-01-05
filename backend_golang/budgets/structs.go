package budgets

type Budget struct {
	BudgetName string `json:"budget_name"`
}

type BudgetItems struct {
	ItemName    string  `json:"item_name"`
	BudgetCost  float32 `json:"budget_cost"`
	Description string  `json:"description"`
	Priority    float32 `json:"priority"`
}

type ManageBudgets struct {
	Budget      Budget        `json:"budget"`
	BudgetItems []BudgetItems `json:"budget_items"`
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
