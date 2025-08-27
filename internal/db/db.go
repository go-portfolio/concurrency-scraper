package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(user, password, host, dbname string, port int) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetURLs() ([]URL, error) {
	rows, err := db.Query("SELECT id, url, created_at FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.ID, &u.URL, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (db *DB) SaveResult(urlID int, content string) error {
	_, err := db.Exec("INSERT INTO results(url_id, content) VALUES($1, $2)", urlID, content)
	return err
}

