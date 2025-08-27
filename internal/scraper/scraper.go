package scraper

import (
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/go-portfolio/concurrency-scraper/internal/models"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/internal/logger"
)

// Scraper — объект скрейпера, принимает интерфейсы
type Scraper struct {
	log    logger.Logger
	client httpclient.Client
	db     db.DB
	pool   worker.Pool
}

// New — конструктор с инъекцией зависимостей
func New(log logger.Logger, client httpclient.Client, database db.DB, pool worker.Pool) *Scraper {
	return &Scraper{
		log:    log,
		client: client,
		db:     database,
		pool:   pool,
	}
}

// Run — основной рабочий цикл
func (s *Scraper) Run(workers int) error {
	var wg sync.WaitGroup

	urls, err := s.db.GetURLs()
	if err != nil {
		return err
	}

	results := make(chan models.ScrapeResult, 10)

	wg.Add(2)

	// Writer goroutine
	go func() {
		defer wg.Done()
		for r := range results {
			resultID, err := s.db.SaveResult(r.URLID, r.Content)
			if err != nil {
				s.log.Error("Ошибка сохранения results: %v", err)
				continue
			}
			if err := s.db.SavePageData(r, resultID); err != nil {
				s.log.Error("Ошибка сохранения pages: %v", err)
			} else {
				s.log.Info("Сохранено: %s (result_id=%d)", r.URL, resultID)
			}
		}
	}()

	// Submitter goroutine
	go func() {
		defer wg.Done()
		s.log.Info("Количество URL для обработки: %d", len(urls))
		for _, u := range urls {
			u := u
			s.pool.Submit(func() {
				body, err := s.client.Fetch(u.URL)
				if err != nil {
					s.log.Error("Ошибка загрузки %s: %v", u.URL, err)
					return
				}

				doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
				if err != nil {
					s.log.Error("Ошибка парсинга %s: %v", u.URL, err)
					return
				}

				title := doc.Find("title").First().Text()
				summary, _ := doc.Find(`meta[name="description"]`).Attr("content")
				if summary == "" {
					summary, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
				}
				text := doc.Find("body").Text()
				wordCount := len(strings.Fields(text))

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

		// дождёмся выполнения всех задач в пуле и закроем results
		s.pool.Close()
		close(results)
	}()

	wg.Wait()
	return nil
}
