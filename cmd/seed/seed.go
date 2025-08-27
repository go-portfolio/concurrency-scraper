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
		"https://golang.org",
		"https://go.dev",
		"https://news.ycombinator.com",                  // Хакерньюс — довольно открытый
		"https://reddit.com/r/golang",                   // Reddit Golang subreddit
		"https://dev.to/t/golang",                       // Dev.to тег Golang
		"https://stackoverflow.com/questions/tagged/go", // СтэкОверфлоу с тегом Go
		"https://pkg.go.dev",                            // Документация Go пакетов
		"https://blog.golang.org",                       // Официальный блог Go
		"https://go101.org",                             // go101.org — бесплатный учебник по Go
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
