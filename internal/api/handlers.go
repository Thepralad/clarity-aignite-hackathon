package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/core"
	"github.com/Thepralad/clarity-aignite-hackathon/internal/db"
)

func HandlerSearch(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := req.URL.Query().Get("query")
	responses, err := core.Crawl(query)
	if err != nil {
		log.Print(err)
	}

	urls := core.ExtractUrls(responses)

	scrapedArticles := core.ScrapeUrls(urls)

	summary, err := core.Summarize(query, &scrapedArticles)
	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to summarize: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the summary to JSON
	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to encode summary to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(summaryJSON)

}

func HandlerArticles(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	res.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	articles, err := db.GetArticles()
	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to get articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the articles to JSON
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to encode articles to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(articlesJSON)
}
