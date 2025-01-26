package budgets

import (
	"database/sql"
	"fmt"
)

func BudgetExists(db *sql.DB, user_id string, budget_name string) bool {
	var res bool
	fmt.Println(budget_name)
	_ = db.QueryRow("SELECT EXISTS(SELECT * FROM Budget WHERE budget_name = ? AND user_id = ?)", budget_name, user_id).Scan(&res)
	return res
}

func ItemExists(db *sql.DB, user_id string, item_name string, budget_name string) bool {
	var res bool
	_ = db.QueryRow("SELECT EXISTS(SELECT * FROM Budget_Items WHERE item_name = ? AND user_id = ? AND budget_name = ?)",
		item_name, user_id, budget_name).Scan(&res)
	return res
}
