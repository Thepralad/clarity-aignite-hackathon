package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/api"
	"github.com/Thepralad/clarity-aignite-hackathon/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: No .env file found, using system environment variables")
	}

	// Fetch config from environment variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	// Initialize the database connection
	err := db.Init(dsn)
	if err != nil {
		fmt.Printf("Database initialization error: %s\n", err)
		return
	}
	http.HandleFunc("/search", api.HandlerSearch)
	http.HandleFunc("/articles", api.HandlerArticles)

	// http.HandleFunc("/news", newsHandler)

	fmt.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server error: %s\n", err)
	}

}

// func searchHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the Search Page!")
// }

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the News Page!")
// }
