package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	
	_ "github.com/go-sql-driver/mysql"
	 "github.com/joho/godotenv"
)

var DB *sql.DB

func InitDatabase() (*sql.DB, error) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env.local"); 
		err != nil {
			fmt.Println("Warning: .env.local not found (non-production)")
		}
	}

	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	dbPort := os.Getenv("MYSQL_PORT")
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, database)
	fmt.Printf("Connecting to DB at %s:%s with user %s\n", dbHost, dbPort, dbUser)

	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)
	DB.SetConnMaxLifetime(10 * time.Minute)

	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("Database connected successfully!")
	return DB, nil
}

func CloseDatabase() error {
	if DB != nil {
		fmt.Println("Closing database connection...")
		return DB.Close()
	}
	return nil
}