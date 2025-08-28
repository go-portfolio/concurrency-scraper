package main

import (
	"log"
	"runtime"

	"github.com/go-portfolio/concurrency-scraper/internal/config"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/es"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Загружаем конфигурацию приложения из файла
	appCfg := loadConfig()

	// Инициализируем логгер
	logr := logger.NewStdLogger()

	// Инициализируем подключение к базе данных
	sqlDB := mustInitDB(appCfg)

	// Инициализируем клиент Elasticsearch (с проверкой подключения)
	esClient := mustInitElastic()

	// Создаём HTTP-клиент
	httpc := httpclient.New()

	// Запускаем пул воркеров с указанным количеством потоков
	pool := worker.NewPool(appCfg.Workers)

	// Инициализируем скрапер с зависимостями (логгер, http, БД, пул)
	s := scraper.New(logr, httpc, sqlDB, pool, esClient)

	// Запускаем основной процесс скрапинга
	if err := s.Run(appCfg.Workers); err != nil {
		logr.Error("run error: %v", err)
	}
}

// loadConfig загружает конфигурацию из файла, определяя путь до текущего файла
func loadConfig() *config.Config {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("не удалось определить путь до конфигурации")
	}
	return config.Load(filename)
}

// mustInitDB инициализирует подключение к SQL-базе данных, паникует при ошибке
func mustInitDB(cfg *config.Config) *db.SQLDB {
	sqlDB, err := db.NewSQLDB(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

// mustInitElastic инициализирует клиент Elasticsearch и проверяет его доступность
func mustInitElastic() es.Client {
	// Конфигурация ES клиента
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username: "elastic",
		// Password: "yourpassword",
	}

	// Проверяем доступность ES через стандартный клиент
	esCheck, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Ошибка при создании клиента Elasticsearch: %s", err)
	}
	res, err := esCheck.Info()
	if err != nil {
		log.Fatalf("Ошибка при получении информации от Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	// Создаем наш клиент, реализующий интерфейс es.Client
	esClient, err := es.New([]string{"http://localhost:9200"}, "pages_index")
	if err != nil {
		log.Fatalf("Ошибка подключения к Elasticsearch: %v", err)
	}

	return esClient
}