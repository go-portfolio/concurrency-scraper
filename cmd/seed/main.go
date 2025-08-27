package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-portfolio/concurrency-scraper/internal/config"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Не удалось определить путь")
	}

	cfg := config.Load(filename)

	database, err := db.NewSQLDB(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal(err)
	}

	urls := []string{
		"https://golang.org",
		"https://go.dev",
		"https://news.ycombinator.com",
		"https://reddit.com/r/golang",
		"https://dev.to/t/golang",
		"https://stackoverflow.com/questions/tagged/go",
		"https://pkg.go.dev",
		"https://blog.golang.org",
		"https://go101.org",
	}

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
