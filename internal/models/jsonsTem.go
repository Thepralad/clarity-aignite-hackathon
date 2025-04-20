package models

type NewsResponse struct {
	SearchParameters struct {
		Query  string `json:"q"`
		Type   string `json:"type"`
		Engine string `json:"engine"`
	} `json:"searchParameters"`
	News []NewsItem `json:"news"`
}

type NewsItem struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snippet"`
	Date     string `json:"date"`
	Source   string `json:"source"`
	Position int    `json:"position"`
}

type ScrapedArticle struct {
	URL     string   `json:"url"`
	Title   string   `json:"title"`
	Content []string `json:"content"`
}

type SummarizedArticle struct {
	Title           string   `json:"title"`
	Para            []string `json:"para"`
	Points          []string `json:"points"`
	SourcesURL      []string `json:"sources_url"`
	RelatedSearches []string `json:"relatedsearch"`
}

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ImgURL   string `json:"img_url"`
	Category string `json:"category"`
}
