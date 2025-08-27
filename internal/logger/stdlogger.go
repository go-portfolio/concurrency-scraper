package logger

import (
	"log"
	"time"
)

// Простая реализация Logger — для CLI/локального запуска
type StdLogger struct{}

func NewStdLogger() Logger {
	return &StdLogger{}
}

func (l *StdLogger) Info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *StdLogger) Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (l *StdLogger) Debug(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format(time.RFC3339)}, args...)
	log.Printf("[DEBUG %s] "+format, args...)
}

