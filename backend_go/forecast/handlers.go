package forecast

import (
	"encoding/json"
	"golang/auth"
	"golang/budgets"
	"golang/database"
	"net/http"

	"github.com/gorilla/mux"
)

func (forecast *ForecastHandler) ForecastTransactions(w http.ResponseWriter, r *http.Request) {

	// extracts user_id from jwt (performed in jwt middleware)
	user_id, ok := r.Context().Value(auth.UserIDkey).(string)
	if !ok {
		http.Error(w, "Cannot find user id in context", http.StatusUnauthorized)
		return
	}

	var vars = mux.Vars(r)
	// get item name, budget name and months to forecast from url
	budget_name := vars["budget_name"]
	// months := vars["months"]

	// check if budget item exists
	if !budgets.BudgetExists(database.DB, user_id, budget_name) {
		http.Error(w, "Budget does not exist, please try again", http.StatusNotFound)
	} else {
		// get all items and transactions related to a budget
		res := GetBudgetData(database.DB, user_id, budget_name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)

		// resString, err := json.Marshal(res)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // request the python forecasting api
		// req, err := http.NewRequest("POST", "http://0.0.0.0:5000/forecast?months="+months,
		// 	bytes.NewBuffer(resString))
		// req.Header.Set("Content-Type", "application/json")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// client := &http.Client{}
		// response, err := client.Do(req)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer response.Body.Close()

		// // read response from request
		// body, err := io.ReadAll(response.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // output results to the user
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(response.StatusCode)
		// w.Write(body)
	}
}
