package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

// func (user *User) UserExists(db *sql.DB) (bool, error) {
// 	rows, err := db.Query("SELECT * FROM user WHERE email = ?", user.Email)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// }

func (account *Account) ValidateUserAndPassword(db *sql.DB) (bool, error) {
	var db_hash string
	var db_salt []byte
	err := db.QueryRow("SELECT password, user_id, salt FROM User WHERE email = ?", account.User.Email).Scan(&db_hash, &account.UserID, &db_salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("the email or password provided isn't correct, please try again or create a new account")
		}
		return false, err
	}
	account.User.HashPassword(db_salt)
	if account.User.Password == db_hash {
		return true, nil
	}
	return false, nil
}

func (account *Account) ValidateSecurityQuestions(db *sql.DB) (bool, error) {
	rows, err := db.Query("SELECT * FROM Security_Questions WHERE user_id = ?", account.UserID)
	defer rows.Close()

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	// db connection

	config := mysql.Config{
		User:                 "root",
		Passwd:               "Lemonadetv2027!?",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "pff",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// db connection

	// login

	var account Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}
	// check if user details are present in the db
	res, err := account.ValidateUserAndPassword(db)
	if err != nil {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message:    err.Error(),
			StatusCode: 400,
		}
		// builds json response
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
			return
		}
	}
	if res {
		account.ValidateSecurityQuestions(db)
	}
}
