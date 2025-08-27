package db

import "time"

type URL struct {
	ID        int
	URL       string
	CreatedAt time.Time
}

type Result struct {
	ID        int
	URLID     int
	Content   string
	CreatedAt time.Time
}

// Общая структура данных, передаваемая из Scraper → DB
type ScrapeResult struct {
    URL       string
    URLID     int
    Title     string
    Summary   string
    Language  string
    WordCount int
    FetchedAt time.Time
    Content   string // raw HTML
}

type PageData struct {
	URL       string
	URLID     int
	Title     string
	Summary   string
	Language  string
	WordCount int
	FetchedAt time.Time
}
