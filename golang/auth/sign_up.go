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

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// user methods

func (user User) GenerateUserID() string {
	res, err := uuid.Parse(user.Username + user.Forename + user.Surname)
	if err != nil {
		log.Fatal(err)
	}
	return res.String()
}

func (user User) StoreUserInformation(sq Security_Questions) {
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
}

func (user User) ValidateSignUp(db *sql.DB) (bool, error) {
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

func HashAndSaltPassword(password string, salt []byte) string {
	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
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

// main sign up process

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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

	var account Account
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

}

// main sign up process