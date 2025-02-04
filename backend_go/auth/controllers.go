package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (account *Account) ValidateUserAndPassword(db *sql.DB) (bool, error) {
	var db_hash string
	var db_salt []byte
	// check if user exists in db
	err := db.QueryRow("SELECT password, user_id, salt FROM User WHERE email = ?",
		account.User.Email).
		Scan(&db_hash, &account.UserID, &db_salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf(`the email or password provided isn't 
			correct, please try again or create a new account`)
		}
		return false, err
	}
	// check if password matches with  password in db
	account.User.HashPassword(db_salt)
	if account.User.Password == db_hash {
		return true, nil
	}
	return false, nil
}

func (account *Account) ValidateSecurityQuestions(db *sql.DB) (bool, error) {
	// queries all security questions by the user in db
	rows, err := db.Query(`SELECT question, answer FROM Security_Questions 
	WHERE user_id = ?`, account.UserID)
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

	if len(account.Security_Questions) < len(sq) {
		return false, nil
	}

	// makes a variable to check for matched questions and answers
	correct := 0
	for _, i := range sq {
		for _, j := range account.Security_Questions {
			// checks if each question and answer match
			if i.Question == j.Question &&
				strings.TrimSpace(strings.ToLower(i.Answer)) ==
					strings.TrimSpace(strings.ToLower(j.Answer)) {
				// increment by one if match
				correct++
				break
			}
		}
	}

	return correct == len(sq), nil
}

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
	// check if password has lowercase, uppercase, digits, and special characters
	check_password :=
		regexp.MustCompile("[a-z]").MatchString(user.Password) &&
			regexp.MustCompile("[A-Z]").MatchString(user.Password) &&
			regexp.MustCompile("[0-9]").MatchString(user.Password) &&
			regexp.MustCompile("[^a-zA-Z0-9]").MatchString(user.Password)

	return check_password
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

var key = []byte("pfftesting")

func (account *Account) GenerateJWT() (string, error) {
	// set JWT expiration date
	expiration_time := time.Now().Add(720 * time.Hour)

	// set up information that is stored in the JWT
	claims := &Claims{
		UserID: account.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration_time),
		},
	}
	// encode claims and sign with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func ValidateJWT(token_string string) (*Claims, error) {
	claims := &Claims{}
	// decode jwt to get user_id
	token, err := jwt.ParseWithClaims(token_string, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect signing method")
			}
			return key, nil
		})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		// check if token has expired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		}
		return nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

type ctxkey string

const UserIDkey ctxkey = "user_id"

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (user *User) HashPassword(salt []byte) {
	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(user.Password))
	user.Password = base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
}
