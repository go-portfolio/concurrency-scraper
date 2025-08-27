package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Определяем путь к текущему исполняемому файлу
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)

	// Путь к .env рядом с исполняемым файлом
	envPath := filepath.Join(exeDir, "../../.env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	// Создаём логгер для вывода сообщений
	logr := logger.New()

	// Читаем переменные окружения для подключения к базе данных
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")
	workersStr := os.Getenv("WORKERS") // количество воркеров для скрейпера

	// Проверяем, что все обязательные переменные окружения установлены
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" || dbPortStr == "" {
		log.Fatal("Не все переменные окружения для БД установлены")
	}

	// Конвертируем порт БД из строки в число
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Неверный порт БД: %v", err)
	}

	// Устанавливаем количество воркеров для скрейпера
	workers := 5 // значение по умолчанию
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil {
			workers = w
		}
	}

	// Подключаемся к базе данных
	database, err := db.New(dbUser, dbPassword, dbHost, dbName, dbPort)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём экземпляр Scraper с логгером и подключением к БД
	s := scraper.New(logr, database)

	// Запускаем скрейпер с указанным количеством воркеров
	if err := s.Run(workers); err != nil {
		log.Fatal(err)
	}
}
