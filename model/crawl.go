package model

type CrawlInput struct {
	Url string `json:"url"`
}

type Crawl struct {
	CurrentPage int    `json:"current_page"`
	URL         string `json:"url"`
}
