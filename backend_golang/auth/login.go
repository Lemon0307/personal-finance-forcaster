package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/database"
	"log"
	"net/http"
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
	// check if user exists in db
	err := db.QueryRow("SELECT password, user_id, salt FROM User WHERE email = ?", account.User.Email).
		Scan(&db_hash, &account.UserID, &db_salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("the email or password provided isn't correct, please try again or create a new account")
		}
		return false, err
	}
	// check if password matches with password in db
	account.User.HashPassword(db_salt)
	if account.User.Password == db_hash {
		return true, nil
	}
	return false, nil
}

func (account *Account) ValidateSecurityQuestions(db *sql.DB) (bool, error) {
	// queries all security questions by the user in db
	rows, err := db.Query("SELECT question, answer FROM Security_Questions WHERE user_id = ?", account.UserID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var sq []Security_Questions
	// convert sql rows into array of struct
	for rows.Next() {
		var question string
		var answer string
		if err := rows.Scan(&question, &answer); err != nil {
			log.Fatal(err)
		}
		sq = append(sq, Security_Questions{question, answer})
	}
	// check if answers to questions match
	for i := 0; i < len(sq); i++ {
		if account.Security_Questions[i].Answer != sq[i].Answer {
			return false, nil
		}
	}
	return true, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	// login

	var account Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}
	// check if user details are present in the db
	user_password_ok, err := account.ValidateUserAndPassword(database.DB)
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
	if user_password_ok {
		// checks if security questions match
		security_questions_ok, _ := account.ValidateSecurityQuestions(database.DB)
		if security_questions_ok {
			// generates jwt token for the user to authenticate
			token, err := account.GenerateJWT()
			if err != nil {
				log.Fatal(err)
			}
			// return message
			w.Header().Set("Content-Type", "application/json")
			response := ResponseJWT{
				Message:    "Successfully logged in!",
				Token:      token,
				StatusCode: 201,
			}
			// builds json response
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
				return
			}
		}
	}
}
