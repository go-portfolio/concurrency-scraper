package es

import "github.com/go-portfolio/concurrency-scraper/internal/models"

// Client — интерфейс для работы с Elasticsearch.
// Позволяет индексировать страницы и закрывать соединение (если необходимо).
type Client interface {
	// IndexPage сохраняет или обновляет документ в Elasticsearch
	// result — структура ScrapeResult, содержащая данные страницы для индексации.
	// Возвращает ошибку в случае проблем с индексацией.
	IndexPage(result models.ScrapeResult) error

	// Close закрывает соединение с Elasticsearch.
	// На практике большинство клиентов ES не требуют явного закрытия, но метод
	// оставлен для совместимости с интерфейсом и моков для тестирования.
	Close() error
}
