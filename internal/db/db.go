package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL драйвер для database/sql
)

// DB — обёртка над стандартной sql.DB для удобной работы с БД
type DB struct {
	*sql.DB
}

// New создаёт новое подключение к базе данных PostgreSQL
// Принимает пользователя, пароль, хост, имя базы и порт
func New(user, password, host, dbname string, port int) (*DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	// Открываем подключение к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Возвращаем обёртку DB
	return &DB{db}, nil
}

// GetURLs возвращает список URL из таблицы urls
func (db *DB) GetURLs() ([]URL, error) {
	// Выполняем SQL-запрос
	rows, err := db.Query("SELECT id, url, created_at FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // закрываем после обработки

	var urls []URL
	for rows.Next() {
		var u URL
		// Сканируем данные строки в структуру URL
		if err := rows.Scan(&u.ID, &u.URL, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

// SaveResult сохраняет результат скрейпинга для конкретного URL в таблицу results
func (db *DB) SaveResult(urlID int, content string) error {
	_, err := db.Exec("INSERT INTO results(url_id, content) VALUES($1, $2)", urlID, content)
	return err
}
