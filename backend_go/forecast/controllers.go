package forecast

import (
	"database/sql"
	"golang/database"
	"log"
)

func GetTransactions(db *sql.DB, user_id, item_name, budget_name string) []TotalTransactions {
	// find month, year and total amount of transactions related to item
	rows, err := database.DB.Query(`
	SELECT t.month, t.year, SUM(t.amount) AS total_amount 
	FROM
		Transactions t
	JOIN 
		Monthly_Costs mc
	ON 
		t.user_id = mc.user_id AND t.month = mc.month AND t.year = mc.year
	WHERE 
		mc.item_name = t.item_name
		AND
		mc.item_name = ?
		AND
		t.budget_name = ?
		AND
		t.user_id = ?
	GROUP BY
		t.month, t.year
	ORDER BY
		t.year, t.month;
	`, item_name, budget_name, user_id)
	if err != nil {
		log.Fatal(err)
	}

	var res []TotalTransactions

	// append all transactions into res
	for rows.Next() {
		var transaction TotalTransactions
		err := rows.Scan(&transaction.Month, &transaction.Year, &transaction.TotalAmount)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, transaction)
	}

	return res
}
