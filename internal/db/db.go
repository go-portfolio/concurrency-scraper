package db

import "github.com/go-portfolio/concurrency-scraper/internal/models"

// DB — интерфейс для слоя доступа к данным (для тестирования)
type DB interface {
	GetURLs() ([]models.URL, error)
	SaveResult(urlID int, content string) (int, error)
	SavePageData(r models.ScrapeResult, resultID int) error
	Exec(query string, args ...interface{}) (interface{}, error) // упрощённо для seed
}
