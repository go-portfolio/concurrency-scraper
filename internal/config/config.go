package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config хранит все переменные окружения проекта
type Config struct {
	DBUser  string
	DBPass  string
	DBHost  string
	DBName  string
	DBPort  int
	Workers int
}

// Load загружает .env файл и возвращает структуру Config
func Load(file string) *Config {
	fmt.Println(file)
	dir := filepath.Dir(file)
	envPath := ""
	found := false
	
	for {
		envPath = filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// Файл найден, загружаем
			if err := godotenv.Load(envPath); err != nil {
				log.Fatal("Ошибка загрузки .env файла:", err)
			}
			found = true
			break
		}

		// Поднимаемся на уровень выше
		parent := filepath.Dir(dir)
		if parent == dir {
			// Достигли корня файловой системы
			break
		}
		dir = parent
	}

	if !found {
		log.Fatal(".env файл не найден")
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	cfg := &Config{}

	// Читаем обязательные переменные
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPass = os.Getenv("DB_PASSWORD")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBName = os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")
	workersStr := os.Getenv("WORKERS")

	if cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBHost == "" || cfg.DBName == "" || dbPortStr == "" {
		log.Fatal("Не все переменные окружения для БД установлены")
	}

	// Конвертируем порт
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Неверный порт БД: %v", err)
	}
	cfg.DBPort = dbPort

	// Количество воркеров
	cfg.Workers = 5
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil {
			cfg.Workers = w
		}
	}

	return cfg
}

// Helper для форматирования строки подключения к PostgreSQL
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
