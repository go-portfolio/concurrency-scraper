package main

import (
	"log"
	"os"
	"strconv"

	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	logr := logger.New()

	// Чтение переменных окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")
	workersStr := os.Getenv("WORKERS")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" || dbPortStr == "" {
		log.Fatal("Не все переменные окружения для БД установлены")
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Неверный порт БД: %v", err)
	}

	workers := 5 // значение по умолчанию
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil {
			workers = w
		}
	}

	// Подключение к базе
	database, err := db.New(dbUser, dbPassword, dbHost, dbName, dbPort)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём скрейпер
	s := scraper.New(logr, database)

	// Запуск скрейпера с указанным количеством воркеров
	if err := s.Run(workers); err != nil {
		log.Fatal(err)
	}
}
