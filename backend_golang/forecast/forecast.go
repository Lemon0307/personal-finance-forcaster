package forecast

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetTransactions(db *sql.DB, user_id, item_name string) []TotalTransactions {
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
	var err error
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID
	var vars = mux.Vars(r)
	item_name := vars["item_name"]
	months := vars["months"]

	if !budgets.ItemExists(database.DB, user_id, item_name) {
		http.Error(w, "Budget item does not exist, please try again", http.StatusNotFound)
	} else {
		res := GetTransactions(database.DB, user_id, item_name)

		resString, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", "http://localhost:5000/forecast?months="+months, bytes.NewBuffer([]byte(resString)))
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

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		w.Write(body)
	}

}

func (forecast *ForecastHandler) RecommendBudget(w http.ResponseWriter, r *http.Request) {
	var err error
	// check if token is valid or expired
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Invalid token, please log in again", http.StatusUnauthorized)
		return
	}
	// get user id from jwt
	user_id := claims.UserID
	var vars = mux.Vars(r)
	item_name := vars["item_name"]

	if !budgets.ItemExists(database.DB, user_id, item_name) {
		http.Error(w, "Budget item does not exist, please try again", http.StatusNotFound)
	} else {
		res := GetTransactions(database.DB, user_id, item_name)

		resString, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", "http://localhost:5000/recommend_budget", bytes.NewBuffer([]byte(resString)))
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

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		w.Write(body)
	}
}

func ForecastRoutes(router *mux.Router, forecastService ForecastService) {
	router.HandleFunc("/forecast/forecast_transactions/{months}/{item_name}", forecastService.ForecastTransactions).Methods("POST")
	router.HandleFunc("/forecast/recommend_budget/{item_name}", forecastService.RecommendBudget).Methods("POST")
}
