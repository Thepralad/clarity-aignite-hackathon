package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"sync"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/models"
	"github.com/gocolly/colly/v2"
)

func Crawl(query string) (models.NewsResponse, error) {
	url := "https://google.serper.dev/news"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"q":"%s"}`, query))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return models.NewsResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("X-API-KEY", "5f62ba3758a4a3b4ab386bb154c98b921a58155c")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return models.NewsResponse{}, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.NewsResponse{}, fmt.Errorf("error reading response: %w", err)
	}

	var newsResponse models.NewsResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return models.NewsResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return newsResponse, nil
}

func ExtractUrls(res models.NewsResponse) []string {
	var urls []string

	// Iterate through each NewsItem in the News slice
	for _, newsItem := range res.News {
		// Add the Link (URL) to our urls slice
		if newsItem.Link != "" {
			urls = append(urls, newsItem.Link)
		}
	}

	return urls[0:5]
}

func ScrapeUrls(urls []string) []models.ScrapedArticle {
	articles := make([]models.ScrapedArticle, 0)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Create a new collector with some basic settings
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Process each URL concurrently
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Create article instance for this URL
			article := models.ScrapedArticle{
				URL:     url,
				Content: make([]string, 0),
			}

			// Clone collector for concurrent scraping
			clone := c.Clone()

			// Extract title (h1)
			clone.OnHTML("h1", func(e *colly.HTMLElement) {
				if article.Title == "" {
					article.Title = strings.TrimSpace(e.Text)
				}
			})

			// Extract paragraphs
			clone.OnHTML("p", func(e *colly.HTMLElement) {
				text := strings.TrimSpace(e.Text)
				if text != "" {
					article.Content = append(article.Content, text)
				}
			})

			// Visit the URL
			err := clone.Visit(url)
			if err == nil {
				// Only add articles that were successfully scraped
				mutex.Lock()
				articles = append(articles, article)
				mutex.Unlock()
			}

		}(url)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return articles
}
