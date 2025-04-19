package main

import (
	"fmt"
	"log"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/core"
)

func main() {
	// http.HandleFunc("/search", searchHandler)

	// http.HandleFunc("/news", newsHandler)

	// fmt.Println("Server starting on port 8080...")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	fmt.Printf("Server error: %s\n", err)
	// }

	responses, err := core.Crawl("america and china tariff senario, who is right?")
	if err != nil {
		log.Print(err)
	}
	fmt.Println(responses)

	urls := core.ExtractUrls(responses)

	// fmt.Println("\nFound URLs:")
	// fmt.Println("-------------------")
	// for i, url := range urls {
	// 	fmt.Printf("%2d. %s\n", i+1, url)
	// }
	// fmt.Println("-------------------")

	scrapedArticles := core.ScrapeUrls(urls)

	fmt.Println(scrapedArticles)
}

// func searchHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the Search Page!")
// }

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the News Page!")
// }
