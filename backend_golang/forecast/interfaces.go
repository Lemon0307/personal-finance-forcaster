package forecast

import "net/http"

type ForecastService interface {
	ForecastTransactions(w http.ResponseWriter, r *http.Request)
}
