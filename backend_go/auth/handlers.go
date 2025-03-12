package auth

import (
	"encoding/json"
	"golang/database"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (auth *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {

	var err error

	var account Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}

	// check if user details are present in the db
	if !account.ValidateUserAndPassword(database.DB) {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `The password you entered is incorrect, please try again`,
			http.StatusUnauthorized)
	} else if !account.SecurityQuestionsValid(database.DB) {
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `Not enough security questions or security questions are invalid, please try again`,
			http.StatusUnauthorized)
	} else {
		// checks if security questions match
		token, err := account.GenerateJWT()
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		}
		// return message
		w.Header().Set("Content-Type", "application/json")
		// build response struct
		response := struct {
			Message, Token string
			StatusCode     int
		}{
			"Successfully logged in!",
			token,
			201,
		}
		// builds json response
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "JSON response could not be encoded",
				http.StatusInternalServerError)
			return
		}
	}
}

func (auth *AuthenticationHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var err error

	var account *Account
	// parse json into account struct
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if UserExists(account.User, database.DB) { // check if data provided exists in the user table
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "A user with this email already has an account, please try again",
			http.StatusConflict)
	} else if !account.User.ValidPassword() { // check if password matches the regex
		// return error message
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `Password must contain lowercase, uppercase, digits, 
		or special characters.`,
			http.StatusUnauthorized)
	} else if account.User.Password != account.User.ConfirmPassword {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Password does not match confirm password",
			http.StatusUnauthorized)
	} else {
		salt := GenerateSalt(16)
		account.User.Salt = salt
		account.User.HashPassword(salt)
		account.UserID = GenerateUserID()
		// add details into the user table
		_, err := database.DB.Exec(`INSERT INTO user (user_id, 
		username, email, password, salt, forename, surname, dob, 
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
			account.User.CurrentBalance)
		// check if query returns errors
		if err != nil {
			http.Error(w, "Could not insert user data into database",
				http.StatusInternalServerError)
		}
		// adds all security questions into security questions table
		for i := 0; i < len(account.Security_Questions); i++ {
			_, err := database.DB.Exec(`INSERT INTO security_questions
				(user_id, question, answer) VALUES (?, ?, ?)`,
				account.UserID,
				account.Security_Questions[i].Question,
				account.Security_Questions[i].Answer)
			// check if query returns errors
			if err != nil {
				http.Error(w, "Could not insert security questions into database",
					http.StatusInternalServerError)
				log.Fatal(err)
			}
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
}

// upgrades the http protocol to websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

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

	// upgrade the connection to websocket
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

	// send back json data with current balance
	err = connection.WriteJSON(map[string]interface{}{
		"current_balance": current_balance,
	})
	if err != nil {
		log.Fatal(err)
	}
}
