package auth

import (
	"database/sql"
	"encoding/json"
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
	token, err := account.GenerateJWT()
	if err != nil {
		log.Println("can't generate JWT: ", err)
		return
	}
}
