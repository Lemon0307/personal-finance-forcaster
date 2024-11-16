package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	Forename        string  `json:"forename"`
	Surname         string  `json:"surname"`
	DOB             Date    `json:"dob"`
	Address         string  `json:"address"`
	CurrentBalance  float32 `json:"current_balance"`
}

type Security_Questions struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Account struct {
	UserID             string
	User               User               `json:"user"`
	Security_Questions []Security_Questions `json:"security_questions"`
}

type Date struct {
	time.Time
}

type Response struct {
	Message    string
	StatusCode int
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
