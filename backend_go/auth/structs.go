package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)
type Account struct {
	UserID             string               `json:"user_id"`
	User               User                 `json:"user"`
	Security_Questions []Security_Questions `json:"security_questions"`
}
type User struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Salt            []byte
	ConfirmPassword string  `json:"confirm_password"`
	DOB             Date    `json:"dob"`
	CurrentBalance  float32 `json:"current_balance"`
}

type Security_Questions struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
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

type AuthenticationHandler struct{}
