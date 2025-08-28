package db

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // библиотека для мокирования базы данных
	"github.com/go-portfolio/concurrency-scraper/internal/models"
	"github.com/stretchr/testify/assert" // удобные ассёрты для тестов
)

func TestGetURLs(t *testing.T) {
	// Создаём "мок" SQL-базы
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Подготавливаем фейковые строки результата
	rows := sqlmock.NewRows([]string{"id", "url", "created_at"}).
		AddRow(1, "https://example.com", time.Now()).
		AddRow(2, "https://golang.org", time.Now())

	// Ожидаем выполнение запроса SELECT и возвращаем подготовленные строки
	mock.ExpectQuery("SELECT id, url, created_at FROM urls").
		WillReturnRows(rows)

	// Подключаем мок вместо настоящей базы
	database := &SQLDB{DB: sqlDB}

	// Вызываем тестируемый метод
	urls, err := database.GetURLs()

	// Проверяем, что ошибок нет и результат соответствует ожиданию
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
	assert.Equal(t, "https://example.com", urls[0].URL)
}

func TestSaveResult(t *testing.T) {
	// Создаём мок базы
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Ожидаем INSERT в таблицу results с конкретными аргументами
	// Возвращаем ID = 42
	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO results (url_id, content, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id
    `)).
		WithArgs(1, "test content").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	database := &SQLDB{DB: sqlDB}

	// Вызываем метод и проверяем результат
	id, err := database.SaveResult(1, "test content")
	assert.NoError(t, err)
	assert.Equal(t, 42, id)
}

func TestSavePageData(t *testing.T) {
	// Создаём мок базы
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Подготавливаем тестовый ScrapeResult
	r := models.ScrapeResult{
		URL:       "https://example.com",
		URLID:     1,
		Title:     "Example",
		Summary:   "Test summary",
		WordCount: 100,
		FetchedAt: time.Now(),
	}

	// Ожидаем INSERT в таблицу pages с аргументами из r
	// (и поведение UPSERT через ON CONFLICT)
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO pages (url, url_id, title, summary, word_count, fetched_at, result_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (url_id) DO UPDATE SET
            title = EXCLUDED.title,
            summary = EXCLUDED.summary,
            word_count = EXCLUDED.word_count,
            fetched_at = EXCLUDED.fetched_at,
            result_id = EXCLUDED.result_id
    `)).
		WithArgs(r.URL, r.URLID, r.Title, r.Summary, r.WordCount, r.FetchedAt, 99).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	database := &SQLDB{DB: sqlDB}

	// Вызываем метод и убеждаемся, что ошибок нет
	err = database.SavePageData(r, 99)
	assert.NoError(t, err)
}
