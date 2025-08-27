package scraper

import (
	"strings"
	"time"

	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"

	"github.com/PuerkitoBio/goquery"
)

// Scraper — главный объект для запуска веб-скрейпера.
// Содержит логгер для вывода сообщений, HTTP-клиент для загрузки страниц и доступ к базе данных.
type Scraper struct {
	log    logger.Logger     // Логгер для вывода ошибок и информации
	client httpclient.Client // HTTP-клиент для загрузки страниц
	db     *db.DB            // Подключение к базе данных для получения URL и сохранения результатов
}

// New создаёт новый экземпляр Scraper с переданным логгером и базой данных.
// HTTP-клиент создаётся автоматически.
func New(log logger.Logger, database *db.DB) *Scraper {
	return &Scraper{
		log:    log,
		client: httpclient.New(), // Создаём новый HTTP-клиент
		db:     database,
	}
}

// Run запускает процесс скрейпинга с указанным количеством воркеров.
// 1. Получает список URL из базы данных.
// 2. Создаёт пул воркеров.
// 3. Отправляет задачи на загрузку страниц.
// 4. Сохраняет результаты в базе данных.
func (s *Scraper) Run(workers int) error {
	// Получаем список URL из базы
	urls, err := s.db.GetURLs()
	if err != nil {
		return err
	}

	// Создаём пул воркеров
	pool := worker.NewPool(workers)

	// Канал для передачи результатов между воркерами и основной горутиной
	results := make(chan struct {
		URL       string
		URLID     int
		Title     string
		Summary   string
		Language  string
		WordCount int
		FetchedAt time.Time
	})

	// Отправляем задачи в пул воркеров
	go func() {
		for _, u := range urls {
			u := u // локальная копия
			pool.Submit(func() {
				body, err := s.client.Fetch(u.URL)
				if err != nil {
					s.log.Error("Ошибка загрузки %s: %v", u.URL, err)
					return
				}

				// Парсим HTML
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

				// Подсчет слов
				text := doc.Find("body").Text()
				wordCount := len(strings.Fields(text))

				results <- db.PageData{
					URL:       u.URL,
					URLID:     u.ID,
					Title:     strings.TrimSpace(title),
					Summary:   strings.TrimSpace(summary),
					WordCount: wordCount,
					FetchedAt: time.Now(),
				}

			})
		}
		pool.Close()
		close(results)
	}()

	return nil
}
