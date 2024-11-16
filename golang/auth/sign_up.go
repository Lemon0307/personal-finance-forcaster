package auth

import (
	// "fmt"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// user methods

func (user *User) GenerateUserID() string {
	res := uuid.New()
	return res.String()
}

func (user *User) ValidateSignUp(db *sql.DB) (bool, error) {
	var res bool
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

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (user *User) HashAndSaltPassword() {
	salt, err := GenerateSalt(16)
	if err != nil {
		log.Fatal(err)
	} else {
		hash := sha256.New()
		hash.Write(salt)
		hash.Write([]byte(user.Password))
		user.Password = base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	}
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
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	res, err := account.User.ValidateSignUp(db)
	if err != nil {
		http.Error(w, "error while validating sign up", http.StatusInternalServerError)
	}
	fmt.Println(res)
	if res {
		if account.User.ValidatePassword() {
			account.User.HashAndSaltPassword()
			account.UserID = account.User.GenerateUserID()
			fmt.Println(account.UserID)
			fmt.Println(account)

			create_user_query, err := db.Exec(`INSERT INTO user (user_id, username, email, password,
			forename, surname, dob, address, current_balance) VALUES 
			(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				account.UserID, account.User.Username, account.User.Email, account.User.Password, account.User.Forename,
				account.User.Surname, account.User.DOB.Time, account.User.Address, account.User.CurrentBalance)
			if err != nil {
				http.Error(w, "Could not insert user data into database", http.StatusInternalServerError)
				log.Fatal(err)
			}

			for i := 0; i < len(account.Security_Questions); i++ {
				create_sq_query, err := db.Exec(`INSERT INTO security_questions
				(user_id, question, answer) VALUES (?, ?, ?)`, account.UserID,
					account.Security_Questions[i].Question, account.Security_Questions[i].Answer)
				if err != nil {
					http.Error(w, "Could not insert sq data into database", http.StatusInternalServerError)
					log.Fatal(err)
				}
				fmt.Println("$1 rows affected in Users Table, $2 rows affected in Security Questions table", create_user_query, create_sq_query)
			}

			w.Header().Set("Content-Type", "application/json")
			response := Response{
				Message:    "Account created!",
				StatusCode: 201,
			}

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			response := Response{
				Message:    "Password does not match confirm password",
				StatusCode: 400,
			}

			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
				return
			}
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message:    "A user with this email already has an account, please try again",
			StatusCode: 400,
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "JSON response could not be encoded", http.StatusInternalServerError)
			return
		}
	}
	fmt.Println(account)
}

// main sign up process
