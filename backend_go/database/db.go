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
    dob DATE NOT NULL,
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
    budget_name VARCHAR(100) NOT NULL,
    user_id CHAR(36) NOT NULL,
    PRIMARY KEY (budget_name, user_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table: Budget Items
CREATE TABLE IF NOT EXISTS Budget_Items (
    item_name VARCHAR(100),
    budget_name VARCHAR(100) NOT NULL, 
    user_id CHAR(36) NOT NULL,
    description TEXT,
    budget_cost DECIMAL(10, 2) NOT NULL,
    priority INT NOT NULL,
    PRIMARY KEY (item_name, budget_name, user_id),
    FOREIGN KEY (budget_name, user_id) REFERENCES Budget(budget_name, user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table: Monthly Costs
CREATE TABLE IF NOT EXISTS Monthly_Costs (
    user_id CHAR(36) NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    item_name VARCHAR(100) NOT NULL,
    budget_name VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (item_name, budget_name, user_id) REFERENCES Budget_Items(item_name, budget_name, user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (user_id, month, year, item_name, budget_name),
    INDEX idx_month_year (month, year, user_id)
);

-- Table: Transactions
CREATE TABLE IF NOT EXISTS Transactions (
    transaction_id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    date DATE NOT NULL,
    item_name VARCHAR(100) NOT NULL,
    budget_name VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id, month, year, item_name, budget_name) REFERENCES Monthly_Costs(user_id, month, year, item_name, budget_name) 
    ON DELETE CASCADE ON UPDATE CASCADE
);
`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}
func InitDB() error {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "pffpassword",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "pff",
		AllowNativePasswords: true,
		Params: map[string]string{
			"multiStatements": "true",
		},
	}

	var err error
	DB, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Connected to the database!")
	err = CreateTable(DB)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Fatal("Failed to close DB connection:", err)
		}
	} else {
		log.Println("No active DB connection to close")
	}
}
