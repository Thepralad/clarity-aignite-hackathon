package main

import (
	"fmt"
	"net/http"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/api"
	"github.com/Thepralad/clarity-aignite-hackathon/internal/db"
)

func main() {
	// Initialize the database connection
	err := db.Init()
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
