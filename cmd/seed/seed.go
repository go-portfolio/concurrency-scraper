package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-portfolio/concurrency-scraper/internal/config"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
)

func main() {
	// Загружаем конфигурацию
	// Определяем путь к текущему исходному файлу
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Не удалось определить путь")
	}

	cfg := config.Load(filename)

	// Подключаемся к базе данных
	database, err := db.New(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Данные для seed
	urls := []string{
		"https://golangweekly.com",
		"https://blog.golang.org",
		"https://news.ycombinator.com",
		"https://github.com/golang/go",
		"https://go.dev/blog",
		"https://dev.to/t/golang",
		"https://medium.com/tag/golang",
		"https://go101.org",
	}

	// Загружаем данные в таблицу urls
	for _, u := range urls {
		_, err := database.Exec("INSERT INTO urls(url) VALUES($1)", u)
		if err != nil {
			fmt.Printf("Ошибка вставки URL %s: %v\n", u, err)
		} else {
			fmt.Printf("URL добавлен: %s\n", u)
		}
	}

	fmt.Println("Seed завершён")
}
