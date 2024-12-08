package budgets

type Budget struct {
	budget_name string `json:"budget_name"`
}

type BudgetItems struct {
	budget_name string  `json:"budget_name"`
	item_name   string  `json:"item_name"`
	budget_cost string  `json:"budget_cost"`
	description string  `json:"description"`
	priority    float32 `json:"priority"`
}

type ManageBudgets struct {
	user_id      string
	budget       Budget        `json:"budget"`
	budget_items []BudgetItems `json:"budget_items"`
}
