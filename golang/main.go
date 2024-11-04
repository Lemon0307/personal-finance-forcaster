package main

import (
	"golang/auth"
	"net/http"

	// "golang/db"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/sign_up", auth.SignUpHandler)
	http.HandleFunc("/login", auth.LoginHandler)

	config := mysql.Config{
		User:   "root",
		Passwd: "Lemonadetv2027!?",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "pff",
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
