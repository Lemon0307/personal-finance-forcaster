package forecast

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTransactions(db *sql.DB, user_id, item_name string) []TotalTransactions {
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
		mc.item_name = ?
		AND
		t.user_id = ?
	GROUP BY
		t.month, t.year
	ORDER BY
		t.year, t.month;
	`, item_name, user_id)
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

func (forecast *ForecastHandler) ForecastTransactions(w http.ResponseWriter, r *http.Request) {

	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var vars = mux.Vars(r)
	// get item name, budget name and months to forecast from url
	item_name := vars["item_name"]
	budget_name := vars["budget_name"]
	months := vars["months"]

	// check if budget item exists
	if !budgets.ItemExists(database.DB, user_id, item_name, budget_name) {
		http.Error(w, "Budget item does not exist, please try again", http.StatusNotFound)
	} else {
		// get all transaction related to item name
		res := GetTransactions(database.DB, user_id, item_name)

		resString, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
		// request the python forecasting api
		req, err := http.NewRequest("POST", "http://localhost:5000/forecast?months="+months,
			bytes.NewBuffer([]byte(resString)))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// read response from request
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// output results to the user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		w.Write(body)
	}
}

func ForecastRoutes(router *mux.Router, forecastService ForecastService) {
	router.HandleFunc("/forecast/forecast_transactions/{months}/{budget_name}/{item_name}",
		forecastService.ForecastTransactions).Methods("POST")
}
