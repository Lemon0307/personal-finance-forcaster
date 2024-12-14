package forecast

type ForecastHandler struct{}

type TotalTransactions struct {
	Month       int
	Year        int
	TotalAmount float64
}

type Response struct {
	ForecastedTransactions []float64
}
