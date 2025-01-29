package auth

import (
	"encoding/json"
	"fmt"
	"golang/database"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (auth *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {

	var err error

	// login

	var account Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account)
	// check if user details are present in the db
	user_password_ok, err := account.ValidateUserAndPassword(database.DB)
	if err != nil {
		// return error message
		http.Error(w, err.Error(), http.StatusNotFound)
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
				http.Error(w, "JSON response could not be encoded",
					http.StatusInternalServerError)
				return
			}
		} else {
			// return error message
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `Not enough security questions or security questions are invalid, please try again`,
				http.StatusUnauthorized)
		}
	}
}

func (auth *AuthenticationHandler) SignUp(w http.ResponseWriter, r *http.Request) {

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
		// check if password matches the regex
		if account.User.ValidatePassword() {
			if account.User.Password == account.User.ConfirmPassword {
				salt, err := GenerateSalt(16)
				account.User.Salt = salt
				if err != nil {
					log.Fatal(err)
				} else {
					account.User.HashPassword(salt)
					account.UserID = GenerateUserID()
					// add details into the user table
					create_user_query, err := database.DB.Exec(`INSERT INTO user (user_id, 
				username, email, password, salt, forename, surname, dob, address, 
				current_balance) VALUES 
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
						http.Error(w, "Could not insert user data into database",
							http.StatusInternalServerError)
						log.Fatal(err)
					}
					// adds all security questions into security questions table
					for i := 0; i < len(account.Security_Questions); i++ {
						create_sq_query, err := database.DB.Exec(`INSERT INTO security_questions
				(user_id, question, answer) VALUES (?, ?, ?)`,
							account.UserID,
							account.Security_Questions[i].Question,
							account.Security_Questions[i].Answer)
						// check if query returns errors
						if err != nil {
							http.Error(w, "Could not insert sq data into database",
								http.StatusInternalServerError)
							log.Fatal(err)
						}
						fmt.Println(`$1 rows affected in Users Table, $2 rows affected in 
					Security Questions table`, create_user_query, create_sq_query)
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
						http.Error(w, "JSON response could not be encoded",
							http.StatusInternalServerError)
						return
					}
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "Password does not match confirm password",
					http.StatusUnauthorized)
			}
		} else {
			// return error message
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `Password must contain lowercase, uppercase, digits, 
			or special characters.`,
				http.StatusUnauthorized)
		}
	} else {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "A user with this email already has an account, please try again",
			http.StatusConflict)
	}
}

// upgrades the http protocol to websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// change this later

func (auth *AuthenticationHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	var err error

	// extract token from the url
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "please provide a token", http.StatusUnauthorized)
		return
	}

	//extract user id from the jwt token
	claims, err := ValidateJWT(token)
	if err != nil {
		log.Fatal(err)
	}
	user_id := claims.UserID

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	var current_balance float32

	// get current balance from the database with user id
	err = database.DB.QueryRow("SELECT current_balance FROM User WHERE user_id = ?", user_id).Scan(&current_balance)
	if err != nil {
		log.Fatal(err)
	}
	err = connection.WriteJSON(map[string]interface{}{
		"current_balance": current_balance,
	})
	if err != nil {
		log.Fatal(err)
	}
}
