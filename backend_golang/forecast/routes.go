package forecast

import "github.com/gorilla/mux"

func ForecastRoutes(router *mux.Router, forecastService ForecastService) {
	router.HandleFunc("/forecast/forecast_transactions/{months}/{budget_name}/{item_name}",
		forecastService.ForecastTransactions).Methods("POST")
}
