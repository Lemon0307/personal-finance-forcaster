package forecast

import (
	"database/sql"
	"log"
)

func GetBudgetData(db *sql.DB, user_id, budget_name string) Items {
	// select all transactions based of a user and a budget
	rows, err := db.Query(`
	SELECT
		DATE_FORMAT(t.date, '%Y-%m') AS month_year,
		SUM(CASE WHEN t.type = 'outflow' THEN t.amount ELSE 0 END) AS total_outflow,
		SUM(CASE WHEN t.type = 'inflow' THEN t.amount ELSE 0 END) AS total_inflow,
		t.item_name,
		bi.budget_cost
	FROM
		Transactions t
	JOIN Budget_Items bi 
		ON t.item_name = bi.item_name 
		AND t.budget_name = bi.budget_name 
		AND t.user_id = bi.user_id
	JOIN Monthly_Costs mc 
		ON t.item_name = mc.item_name 
		AND t.budget_name = mc.budget_name 
		AND t.user_id = mc.user_id 
		AND t.month = mc.month 
		AND t.year = mc.year
	WHERE
		bi.budget_name = ?
		AND t.user_id = ?
	GROUP BY
		DATE_FORMAT(t.date, '%Y-%m'),
		bi.description,
		t.item_name,
		bi.budget_cost
	ORDER BY
		month_year;
	`, budget_name, user_id)
	if err != nil {
		log.Fatal(err)
	}

	// make a hash map to store transactions for items
	itemMap := make(map[string]*Item)

	for rows.Next() {
		var item_name string
		var total_spent TotalTransactions
		var total_earned TotalTransactions
		var transaction_date []uint8
		var budget_cost float64

		// gather data returned from the query
		err := rows.Scan(&transaction_date, &total_spent.Amount, &total_earned.Amount, &item_name, &budget_cost)
		if err != nil {
			log.Fatal(err)
		}
		total_spent.Date = string(transaction_date)

		// make a new key value pair if item doesn't exist
		if _, exists := itemMap[item_name]; !exists {
			itemMap[item_name] = &Item{
				ItemName:    item_name,
				BudgetCost:  budget_cost,
				TotalSpent:  []TotalTransactions{},
				TotalEarned: []TotalTransactions{},
			}
		}

		// add total spent and total earned for an item
		itemMap[item_name].TotalSpent = append(itemMap[item_name].TotalSpent, total_spent)
		itemMap[item_name].TotalEarned = append(itemMap[item_name].TotalEarned, total_earned)
	}

	// turn items hash map into an items array
	var items []Item
	for _, item := range itemMap {
		items = append(items, *item)
	}

	return Items{BudgetName: budget_name, Items: items}
}
