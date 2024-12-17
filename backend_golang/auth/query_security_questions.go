package auth

import (
	"encoding/json"
	"fmt"
	"golang/database"
	"log"
	"net/http"
)

type Details struct {
	Email    string
	Password string
}

func SQHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var details Details
	json.NewDecoder(r.Body).Decode(&details)
	fmt.Print(details)
	var user_id string

	err = database.DB.QueryRow("SELECT user_id FROM User WHERE email = ?",
		details.Email).Scan(&user_id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := database.DB.Query(`SELECT question FROM Security_Questions 
	WHERE user_id = ?`, user_id)
	if err != nil {
		log.Fatal(err)
	}
	var questions []string
	for rows.Next() {
		var question string
		if rows.Scan(&question); err != nil {
			log.Fatal(err)
		}
		questions = append(questions, question)
	}

	response := questions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
