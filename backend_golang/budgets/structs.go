package budgets

type Budget struct {
	BudgetName string `json:"budget_name"`
}

type BudgetItems struct {
	BudgetName  string  `json:"budget_name"`
	ItemName    string  `json:"item_name"`
	BudgetCost  string  `json:"budget_cost"`
	Description string  `json:"description"`
	Priority    float32 `json:"priority"`
}

type ManageBudgets struct {
	UserID      string
	Budget      Budget        `json:"budget"`
	BudgetItems []BudgetItems `json:"budget_items"`
}

type BudgetHandler struct{}
