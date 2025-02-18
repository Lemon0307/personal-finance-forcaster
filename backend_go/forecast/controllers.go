package forecast

import (
	"database/sql"
	"log"
)

func GetBudgetData(db *sql.DB, user_id, budget_name string) Items {
	rows, err := db.Query(`
	SELECT 
		t.item_name, 
		SUM(CASE WHEN t.type = 'outflow' THEN t.amount END), 
		t.month, 
		t.year, 
		bi.budget_cost
	FROM
		Transactions t
	JOIN 
		Monthly_Costs mc 
		ON t.user_id = mc.user_id 
		AND t.month = mc.month 
		AND t.year = mc.year
	JOIN 
		Budget_Items bi 
		ON t.item_name = bi.item_name 
		AND t.budget_name = bi.budget_name 
		AND t.user_id = bi.user_id
	WHERE 
		t.user_id = ? 
		AND t.budget_name = ?
	GROUP BY 
    	t.budget_name, t.item_name, bi.budget_cost, t.month, t.year;
	`, user_id, budget_name)
	if err != nil {
		log.Fatal(err)
	}

	itemMap := make(map[string]*Item)

	for rows.Next() {
		var item_name string
		var transaction TotalTransactions
		var budgetCost float64
		err := rows.Scan(&item_name, &transaction.Amount, &transaction.Month, &transaction.Year, &budgetCost)
		if err != nil {
			log.Fatal(err)
		}

		if _, exists := itemMap[item_name]; !exists {
			itemMap[item_name] = &Item{
				ItemName:          item_name,
				BudgetCost:        budgetCost,
				TotalTransactions: []TotalTransactions{},
			}
		}

		itemMap[item_name].TotalTransactions = append(itemMap[item_name].TotalTransactions, transaction)
	}

	var items []Item
	for _, item := range itemMap {
		items = append(items, *item)
	}

	return Items{BudgetName: budget_name, Items: items}
}
