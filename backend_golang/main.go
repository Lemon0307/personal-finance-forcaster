package main

import (
	"database/sql"
	"fmt"
	"golang/auth"
	"golang/budgets"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/sign_up", auth.SignUpHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/budgets", budgets.BudgetHandler)

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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
