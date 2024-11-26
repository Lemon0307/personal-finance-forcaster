package budgets

type Budget struct {
	BudgetID   string
	UserID     string
	BudgetName string
}

type BudgetFeatures struct {
	BudgetID   string
	FeatureID  string
	CostBudget float32
	Priority   float32
}

type Feature struct {
	FeatureID   string
	Description string
	DefaultCost float32
}
