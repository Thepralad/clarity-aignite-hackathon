package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/core"
)

func HandlerSearch(res http.ResponseWriter, req *http.Request) {

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
