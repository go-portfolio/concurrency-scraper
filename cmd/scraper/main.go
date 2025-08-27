package main

import (
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

// main — точка входа в программу.
// Здесь инициализируем логгер, создаём список URL и запускаем Scraper.
func main() {
	// Создаём новый логгер для вывода информации и ошибок.
	log := logger.New()

	// Список URL, которые будем скрейпить.
	// Можно добавлять свои страницы для теста или демонстрации.
	startURLs := []string{
		"https://example.com",
		"https://golang.org",
		"https://httpbin.org",
	}

	// Создаём новый Scraper, передавая логгер.
	s := scraper.New(log)

	// Запускаем Scraper.
	// Второй параметр — количество воркеров (горутин), которые будут параллельно обрабатывать страницы.
	s.Run(startURLs, 5) // 5 workers
}
