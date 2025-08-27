package scraper

import (
	"concurrency-scraper/internal/httpclient"
	"concurrency-scraper/internal/worker"
	"concurrency-scraper/pkg/logger"
)

type Scraper struct {
	log     logger.Logger
	client  httpclient.Client
}

func New(log logger.Logger) *Scraper {
	return &Scraper{
		log:    log,
		client: httpclient.New(),
	}
}

func (s *Scraper) Run(urls []string, workers int) {
	pool := worker.NewPool(workers)

	results := make(chan string)

	// Отправляем задачи
	go func() {
		for _, url := range urls {
			u := url
			pool.Submit(func() {
				body, err := s.client.Fetch(u)
				if err != nil {
					s.log.Error("Ошибка загрузки %s: %v", u, err)
					return
				}
				results <- body[:80] // кусочек HTML для примера
			})
		}
		pool.Close()
	}()

	// Читаем результаты
	for res := range results {
		s.log.Info("Скачан фрагмент: %s...", res)
	}
}
