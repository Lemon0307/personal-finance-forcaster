package transactions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func GenerateTransactionID() string {
	res := uuid.New()
	return res.String()
}

// switch case to map month in string to map in int
func MonthToInt(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	}
	return 0
}

func TransactionExists(db *sql.DB, transaction_id, user_id string, month, year int,
	item_name, budget_name string) bool {
	var res bool
	// query transactions table for a record that contains user id and transaction id
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM transactions WHERE transaction_id = ?
	AND user_id = ? AND month = ? AND year = ? AND item_name = ? AND budget_name = ?)`,
		transaction_id, user_id, month, year, item_name, budget_name).
		Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (date *Date) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])

	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	date.Time = parsedTime
	return nil
}
