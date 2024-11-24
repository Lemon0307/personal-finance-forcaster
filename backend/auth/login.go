package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

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

	var account *Account
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}
	var res bool
	db.QueryRow("SELECT EXISTS(SELECT * FROM user WHERE user_id = ?)", account.UserID).Scan(&res)
	if !res {
		http.Error(w, "Account does not exists, please try again or create an account", http.StatusNotFound)
	}

	if account.UserID == "" {
		http.Error(w, "Invalid credentials, please try again", http.StatusUnauthorized)
	}

	token, err := account.GenerateJWT()
	if err != nil {
		log.Println("can't generate JWT: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
	fmt.Println("Generated token for user " + account.UserID)
}
