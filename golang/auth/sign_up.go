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
	fmt.Println(res.String())
	return res.String()
}

func (user *User) ValidateSignUp(db *sql.DB) (bool, error) {
	var res bool
	query := "SELECT * FROM user WHERE email = ?"
	err := db.QueryRow(query, user.Email).Scan(&res)
	if err != nil {
		return !res, err
	}
	return true, nil
}

// user methods

// password handling methods

func ConfirmPassword(password string, confirm_password string) bool {
	return password == confirm_password
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

// security questions methods

func GenerateSecurityQID(sq Security_Questions) string {
	res, err := uuid.Parse(sq.Question + sq.Answer)
	if err != nil {
		log.Fatal(err)
	}
	return res.String()
}

// security questions methods

// parsing methods

func (date *Date) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])

	parsedTime, err := time.Parse("02-01-2006", s)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	date.Time = parsedTime
	return nil
}

// parsing methods

// main sign up process

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	// db connection

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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// db connection

	var account Account
	err2 := json.NewDecoder(r.Body).Decode(&account)
	if err2 != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Println(account.User.GenerateUserID())
	account.User.HashAndSaltPassword()
	res, err := account.User.ValidateSignUp(db)
	if err != nil {
		log.Fatal(err)
	} else {
		if res {
			
		}
	}
}

// main sign up process
