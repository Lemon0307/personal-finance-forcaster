package auth

import (
	// "fmt"

	"golang/database"

	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/google/uuid"
)

// user methods

func GenerateUserID() string {
	res := uuid.New()
	return res.String()
}

func (user *User) ValidateSignUp(db *sql.DB) (bool, error) {
	var res bool
	// checks if user exists
	query := "SELECT EXISTS(SELECT * FROM user WHERE email = ?)"
	err := db.QueryRow(query, user.Email).Scan(&res)
	if err != nil {
		return false, err
	}
	return !res, nil
}

// user methods

// password handling methods

func (user User) ValidatePassword() bool {
	var res bool = true
	// rx := "[A-Z]+[a-z]+[0-9]"
	if user.Password != user.ConfirmPassword {
		return false
	}
	// regexp.Match(rx, []byte(user.Password))
	return res
}

// password handling methods

// parsing methods

func (date *Date) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])

	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	date.Time = parsedTime
	return nil
}

// parsing methods

// main sign up process

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var account *Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	fmt.Println(account)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// check if data provided exists in the user table
	res, err := account.User.ValidateSignUp(database.DB)
	if err != nil {
		http.Error(w, "error while validating sign up", http.StatusInternalServerError)
	}

	if res {
		// check if password = confirm password
		if account.User.ValidatePassword() {
			salt, err := GenerateSalt(16)
			account.User.Salt = salt
			if err != nil {
				log.Fatal(err)
			} else {
				account.User.HashPassword(salt)
				account.UserID = GenerateUserID()
				// add details into the user table
				create_user_query, err := database.DB.Exec(`INSERT INTO user (user_id, username, email, password, salt,
					forename, surname, dob, address, current_balance) VALUES 
					(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					account.UserID,
					account.User.Username,
					account.User.Email,
					account.User.Password,
					account.User.Salt,
					account.User.Forename,
					account.User.Surname,
					account.User.DOB.Time,
					account.User.Address,
					account.User.CurrentBalance)
				// check if query returns errors
				if err != nil {
					http.Error(w, "Could not insert user data into database", http.StatusInternalServerError)
					log.Fatal(err)
				}
				// adds all security questions into security questions table
				for i := 0; i < len(account.Security_Questions); i++ {
					create_sq_query, err := database.DB.Exec(`INSERT INTO security_questions
				(user_id, question, answer) VALUES (?, ?, ?)`, account.UserID,
						account.Security_Questions[i].Question, account.Security_Questions[i].Answer)
					// check if query returns errors
					if err != nil {
						http.Error(w, "Could not insert sq data into database", http.StatusInternalServerError)
						log.Fatal(err)
					}
					fmt.Println("$1 rows affected in Users Table, $2 rows affected in Security Questions table", create_user_query, create_sq_query)
				}
				// return message
				w.Header().Set("Content-Type", "application/json")
				response := Response{
					Message:    "Account created!",
					StatusCode: 201,
				}
				// builds json response
				err = json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
					return
				}
			}
		} else {
			// return error message
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Password does not match confirm password", http.StatusUnauthorized)
		}
	} else {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "A user with this email already has an account, please try again", http.StatusConflict)
	}
}

// main sign up process
