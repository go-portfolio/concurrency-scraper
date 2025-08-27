package scraper

import (
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

// Scraper — главный объект для запуска веб-скрейпера.
// Он хранит логгер (для вывода сообщений) и HTTP-клиент (для загрузки страниц).
type Scraper struct {
	log    logger.Logger      // интерфейс логгера
	client httpclient.Client  // HTTP-клиент для скачивания страниц
}

// New создаёт новый Scraper.
// Внутри инициализируется логгер и HTTP-клиент.
func New(log logger.Logger) *Scraper {
	return &Scraper{
		log:    log,
		client: httpclient.New(), // создаём новый клиент
	}
}

// Run запускает скрейпер.
// urls — список адресов для скачивания.
// workers — количество воркеров (горутин), которые будут параллельно обрабатывать задачи.
func (s *Scraper) Run(urls []string, workers int) {
	// Создаём пул воркеров с заданным количеством горутин.
	pool := worker.NewPool(workers)

	// Канал для результатов: сюда будем складывать кусочки HTML.
	results := make(chan string)

	// Запускаем отдельную горутину, которая отправляет задачи в пул.
	go func() {
		for _, url := range urls {
			u := url // создаём локальную копию (важно в замыкании!)
			pool.Submit(func() {
				// Скачиваем страницу
				body, err := s.client.Fetch(u)
				if err != nil {
					// Если ошибка, пишем в лог
					s.log.Error("Ошибка загрузки %s: %v", u, err)
					return
				}
				// Отправляем первые 80 символов HTML в канал результатов
				results <- body[:80]
			})
		}
		// После того как все задачи отправлены, закрываем пул
		pool.Close()
		// И закрываем канал результатов, чтобы "читатель" завершил цикл
		close(results)
	}()

	// Читаем результаты из канала и печатаем их в лог.
	// Цикл закончится автоматически, когда канал будет закрыт.
	for res := range results {
		s.log.Info("Скачан фрагмент: %s...", res)
	}
}
