package scraper

import (
	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

// Scraper — главный объект для запуска веб-скрейпера.
// Содержит логгер для вывода сообщений, HTTP-клиент для загрузки страниц и доступ к базе данных.
type Scraper struct {
	log    logger.Logger   // Логгер для вывода ошибок и информации
	client httpclient.Client // HTTP-клиент для загрузки страниц
	db     *db.DB          // Подключение к базе данных для получения URL и сохранения результатов
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
		urlID   int
		content string
	})

	// Отправляем задачи в пул воркеров
	go func() {
		for _, u := range urls {
			u := u // локальная копия для замыкания
			pool.Submit(func() {
				// Загружаем страницу
				body, err := s.client.Fetch(u.URL)
				if err != nil {
					s.log.Error("Ошибка загрузки %s: %v", u.URL, err)
					return
				}
				// Отправляем результат в канал
				results <- struct {
					urlID   int
					content string
				}{urlID: u.ID, content: body}
			})
		}
		pool.Close()   // Закрываем пул после отправки всех задач
		close(results) // Закрываем канал после завершения всех воркеров
	}()

	// Обрабатываем результаты из канала и сохраняем в базу данных
	for r := range results {
		err := s.db.SaveResult(r.urlID, r.content)
		if err != nil {
			s.log.Error("Ошибка сохранения результата для URL ID %d: %v", r.urlID, err)
		} else {
			s.log.Info("Сохранён результат для URL ID %d", r.urlID)
		}
	}

	return nil
}
