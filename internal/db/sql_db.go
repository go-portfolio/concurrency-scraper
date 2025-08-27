package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/go-portfolio/concurrency-scraper/internal/models"
)

type SQLDB struct {
	*sql.DB
}

func NewSQLDB(user, password, host, dbname string, port int) (*SQLDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &SQLDB{db}, nil
}

// Implement DB interface

func (s *SQLDB) GetURLs() ([]models.URL, error) {
	rows, err := s.Query("SELECT id, url, created_at FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var u models.URL
		if err := rows.Scan(&u.ID, &u.URL, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (s *SQLDB) SaveResult(urlID int, content string) (int, error) {
	var id int
	err := s.QueryRow(`
        INSERT INTO results (url_id, content, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id
    `, urlID, content).Scan(&id)
	return id, err
}

func (s *SQLDB) SavePageData(r models.ScrapeResult, resultID int) error {
	_, err := s.Exec(`
        INSERT INTO pages (url, url_id, title, summary, word_count, fetched_at, result_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (url_id) DO UPDATE SET
            title = EXCLUDED.title,
            summary = EXCLUDED.summary,
            word_count = EXCLUDED.word_count,
            fetched_at = EXCLUDED.fetched_at,
            result_id = EXCLUDED.result_id
    `, r.URL, r.URLID, r.Title, r.Summary, r.WordCount, r.FetchedAt, resultID)
	return err
}

func (s *SQLDB) Exec(query string, args ...interface{}) (interface{}, error) {
	res, err := s.DB.Exec(query, args...)
	return res, err
}
