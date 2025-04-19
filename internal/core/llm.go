package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Thepralad/clarity-aignite-hackathon/internal/models"
)

const (
	API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	API_KEY = "AIzaSyDa7iHi_gkmPPRaJNMN1ZlP5-FjMRme4wg"
)

type ContentPart struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []ContentPart `json:"parts"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

func Summarize(query string, scrapedArticles *[]models.ScrapedArticle) (*models.SummarizedArticle, error) {
	// Start with the instruction part of the prompt
	prompt := "Summarize all the content below, which are the different articles of the news, remove any biases, fact check in just 10-100 words and just be on point(use simple english), generate a succinct title and content for now. and generate in a json format like this title, array of strings(content, to denote paragraph,) and relatedsearch. and only use 2 para if needed, or use only one. And answer this query(as a para in the content field) \n\n" + query

	// Append each article's content to the prompt
	for i, article := range *scrapedArticles {
		prompt += fmt.Sprintf("Article %d:\nTitle: %s\nContent:\n", i+1, article.Title)
		for _, paragraph := range article.Content {
			prompt += paragraph + "\n"
		}
		prompt += "---\n"
	}

	// Prepare the request body
	body := RequestBody{
		Contents: []Content{
			{Parts: []ContentPart{{Text: prompt}}},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest("POST", API_URL+"?key="+API_KEY, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(respBody))
	}

	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	rawText := result.Candidates[0].Content.Parts[0].Text

	// Strip ```json and ending ``` if present
	cleaned := rawText
	if len(cleaned) > 6 && cleaned[:7] == "```json" {
		cleaned = cleaned[7:]
	}
	cleaned = string(bytes.Trim([]byte(cleaned), "`\n"))

	var summary models.SummarizedArticle
	if err := json.Unmarshal([]byte(cleaned), &summary); err != nil {
		return nil, fmt.Errorf("failed to parse summarized article: %v\noriginal: %s", err, cleaned)
	}

	return &summary, nil
}
