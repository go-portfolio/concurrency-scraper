package main

import (
	"log"
	"runtime"

	"github.com/go-portfolio/concurrency-scraper/internal/config"
	"github.com/go-portfolio/concurrency-scraper/internal/db"
	"github.com/go-portfolio/concurrency-scraper/internal/httpclient"
	"github.com/go-portfolio/concurrency-scraper/internal/scraper"
	"github.com/go-portfolio/concurrency-scraper/internal/worker"
	"github.com/go-portfolio/concurrency-scraper/pkg/logger"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Не удалось определить путь")
	}

	cfg := config.Load(filename)

	// логгер
	logr := logger.NewStdLogger()

	// база данных (SQL реализация)
	sqlDB, err := db.NewSQLDB(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal(err)
	}

	// http client
	httpc := httpclient.New()

	// пул воркеров
	pool := worker.NewPool(cfg.Workers)

	// scraper (инъекция зависимостей)
	s := scraper.New(logr, httpc, sqlDB, pool)

	if err := s.Run(cfg.Workers); err != nil {
		logr.Error("run error: %v", err)
	}
}
