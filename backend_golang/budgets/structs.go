package budgets

type Budget struct {
	BudgetID   string
	BudgetName string  `json:"budget_name"`
	BudgetCost float32 `json:"budget_cost"`
	Priority   float32 `json:"priority"`
}

type Feature struct {
	FeatureID   string
	Description string  `json:"description"`
	DefaultCost float32 `json:"default_cost"`
}

type ManageBudgets struct {
	UserID  string
	Budget  Budget    `json:"budget"`
	Feature []Feature `json:"feature"`
}
