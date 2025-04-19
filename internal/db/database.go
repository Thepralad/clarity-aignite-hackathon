package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection pool
var DB *sql.DB

// Init initializes the MySQL database connection
func Init() error {
	// Fetch config from environment variables
	user := "avnadmin"
	pass := os.Getenv("DB_PASS")
	host := "mysql-195ced1-thepralad6-5410.j.aivencloud.com"
	port := "18009"
	name := "defaultdb"

	// Fallbacks or panic if not set (optional safety)
	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		return fmt.Errorf("database config not fully set in environment variables")
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}

	// Set connection pool limits
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to MySQL database successfully")
	return nil
}
