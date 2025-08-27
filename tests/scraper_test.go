package tests

import (
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
	"testing"
)

func TestScraper(t *testing.T) {
	log := logger.New()
	s := scraper.New(log)

	urls := []string{"https://httpbin.org/html"}
	s.Run(urls, 2)
}
