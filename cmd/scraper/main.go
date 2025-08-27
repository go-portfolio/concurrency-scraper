package main

import (
	"log"
	"runtime"

	"github.com/go-portfolio/concurrency-scraper/internal/config"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

func main() {
	// Загружаем конфигурацию
	// Определяем путь к текущему исходному файлу
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Не удалось определить путь")
	}

	cfg := config.Load(filename)

	// Создаём логгер
	logr := logger.New()

	// Подключаемся к базе данных
	database, err := db.New(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём скрейпер
	s := scraper.New(logr, database)

	// Запускаем с указанным количеством воркеров
	if err := s.Run(cfg.Workers); err != nil {
		log.Fatal(err)
	}
}
