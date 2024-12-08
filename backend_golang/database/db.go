package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func CreateTable(db *sql.DB) error {
	const query = `
-- Table: User
CREATE TABLE IF NOT EXISTS User (
    user_id CHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    salt VARBINARY(1024) NOT NULL,
    forename VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    dob DATE NOT NULL,
    address TEXT NOT NULL,
    current_balance DECIMAL(65, 2) NOT NULL
);

-- Table: Security Questions
CREATE TABLE IF NOT EXISTS Security_Questions (
    user_id CHAR(36) NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table: Budget
CREATE TABLE IF NOT EXISTS Budget (
    budget_name VARCHAR(100) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table: Budget Items
CREATE TABLE IF NOT EXISTS Budget_Items (
    item_name VARCHAR(100) PRIMARY KEY,
    budget_name VARCHAR(100) NOT NULL, 
    user_id CHAR(36) NOT NULL,
    description TEXT,
    budget_cost DECIMAL(10, 2) NOT NULL,
    priority INT NOT NULL,
    FOREIGN KEY (budget_name) REFERENCES Budget(budget_name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table: Monthly Costs
CREATE TABLE IF NOT EXISTS Monthly_Costs (
    user_id CHAR(36) NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    item_name VARCHAR(100) NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (item_name) REFERENCES Budget_Items(item_name) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (user_id, month, year),
    INDEX idx_month_year (month, year)
);

-- Table: Transactions
CREATE TABLE IF NOT EXISTS Transactions (
    transaction_id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    transaction_name VARCHAR(100) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (month, year) REFERENCES Monthly_Costs(month, year) ON DELETE CASCADE ON UPDATE CASCADE
);
`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

func InitDB() {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "Lemonadetv2027!?",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "pff",
		AllowNativePasswords: true,
		Params: map[string]string{
			"multiStatements": "true",
		},
	}

	DB, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := CreateTable(DB); err != nil {
		panic(fmt.Sprintf("failed to create tables: %v", err))
	}
	fmt.Println("Connected to the database!")
}
