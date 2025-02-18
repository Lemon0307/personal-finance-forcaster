package forecast

import (
	"database/sql"
	"log"
)

func GetBudgetData(db *sql.DB, user_id, budget_name string) Items {
	rows, err := db.Query(`
	SELECT 
		t.item_name, 
		SUM(CASE WHEN t.type = 'outflow' THEN t.amount ELSE 0 END), 
		SUM(CASE WHEN t.type = 'inflow' THEN t.amount ELSE 0 END),
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
		var total_spent TotalTransactions
		var total_earned TotalTransactions
		var year int
		var month int
		var budget_cost float64
		err := rows.Scan(&item_name, &total_spent.Amount, &total_earned.Amount, &month, &year, &budget_cost)
		if err != nil {
			log.Fatal(err)
		}
		total_spent.Month = month
		total_spent.Year = year
		total_earned.Month = month
		total_earned.Year = year

		if _, exists := itemMap[item_name]; !exists {
			itemMap[item_name] = &Item{
				ItemName:    item_name,
				BudgetCost:  budget_cost,
				TotalSpent:  []TotalTransactions{},
				TotalEarned: []TotalTransactions{},
			}
		}

		itemMap[item_name].TotalSpent = append(itemMap[item_name].TotalSpent, total_spent)
		itemMap[item_name].TotalEarned = append(itemMap[item_name].TotalEarned, total_earned)
	}

	var items []Item
	for _, item := range itemMap {
		items = append(items, *item)
	}

	return Items{BudgetName: budget_name, Items: items}
}
