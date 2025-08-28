package scraper

import (
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/es"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/models"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

// Scraper — основной объект скрейпера, хранит зависимости
type Scraper struct {
	log    logger.Logger     // Логгер для вывода информации
	client httpclient.Client // HTTP-клиент для загрузки страниц
	db     db.DB             // Интерфейс работы с базой данных
	pool   worker.Pool       // Пул воркеров для параллельной обработки URL
	es     es.Client         // Клиент Elasticsearch для индексации данных
}

// New — конструктор Scraper с инъекцией зависимостей
func New(log logger.Logger, client httpclient.Client, database db.DB, pool worker.Pool, esClient es.Client) *Scraper {
	return &Scraper{
		log:    log,
		client: client,
		db:     database,
		pool:   pool,
		es:     esClient,
	}
}

// Run — основной рабочий метод скрейпера
// workers — количество одновременно работающих воркеров (используется в пуле)
func (s *Scraper) Run(workers int) error {
	var wg sync.WaitGroup

	// Получаем список URL из базы данных
	urls, err := s.db.GetURLs()
	if err != nil {
		return err
	}

	// Канал для передачи результатов скрейпа между воркерами и writer goroutine
	results := make(chan models.ScrapeResult, 10)

	wg.Add(2) // Две главные горутины: writer и submitter

	// Writer goroutine — сохраняет результаты в БД и Elasticsearch
	go func() {
		defer wg.Done()
		for r := range results {
			// Сохраняем основной результат в таблицу results
			resultID, err := s.db.SaveResult(r.URLID, r.Content)
			if err != nil {
				s.log.Error("Ошибка сохранения results: %v", err)
				continue
			}

			// Сохраняем дополнительные данные страницы (title, summary, word count)
			if err := s.db.SavePageData(r, resultID); err != nil {
				s.log.Error("Ошибка сохранения pages: %v", err)
			} else {
				s.log.Info("Сохранено: %s (result_id=%d)", r.URL, resultID)
			}

			// Индексация в Elasticsearch
			if s.es != nil {
				if err := s.es.IndexPage(r); err != nil {
					s.log.Error("Ошибка индексации в ES: %v", err)
				} else {
					s.log.Info("Индексировано в ES: %s", r.URL)
				}
			}
		}
	}()

	// Submitter goroutine — отправляет задачи на пул воркеров
	go func() {
		defer wg.Done()
		s.log.Info("Количество URL для обработки: %d", len(urls))
		for _, u := range urls {
			u := u // локальная копия для замыкания
			s.pool.Submit(func() {
				// Загружаем страницу
				body, err := s.client.Fetch(u.URL)
				if err != nil {
					s.log.Error("Ошибка загрузки %s: %v", u.URL, err)
					return
				}

				// Парсим HTML с помощью goquery
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
				if err != nil {
					s.log.Error("Ошибка парсинга %s: %v", u.URL, err)
					return
				}

				// Извлекаем заголовок страницы
				title := doc.Find("title").First().Text()

				// Извлекаем описание страницы (meta description)
				summary, _ := doc.Find(`meta[name="description"]`).Attr("content")
				if summary == "" {
					summary, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
				}

				// Подсчитываем количество слов в тексте страницы
				text := doc.Find("body").Text()
				wordCount := len(strings.Fields(text))

				// Отправляем результат в канал для writer goroutine
				results <- models.ScrapeResult{
					URL:       u.URL,
					URLID:     u.ID,
					Title:     strings.TrimSpace(title),
					Summary:   strings.TrimSpace(summary),
					WordCount: wordCount,
					FetchedAt: time.Now(),
					Content:   body,
				}
			})
		}

		// Закрываем пул и канал результатов после обработки всех задач
		s.pool.Close()
		close(results)
	}()

	wg.Wait() // Ожидаем завершения обеих горутин
	return nil
}
