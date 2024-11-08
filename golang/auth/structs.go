package auth

import (
	"time"
)

type User struct {
	UserID         string
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Forename       string    `json:"forename"`
	Surname        string    `json:"surname"`
	DOB            time.Time `json:"dob"`
	Address        string    `json:"address"`
	CurrentBalance float32   `json:"current_balance"`
}

type Security_Questions struct {
	UserID   string
	Question string
	Answer   string
}

type Account struct {
	User User `json:"user"`
	Security_Questions Security_Questions `json:"security_questions"`
}

var config := mysql.Config{
	User:   "root",
	Passwd: "Lemonadetv2027!?",
	Net:    "tcp",
	Addr:   "localhost:3306",
	DBName: "pff",
}