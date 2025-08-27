package logger

import (
	"log"
	"os"
)

// Logger — интерфейс логгера с двумя уровнями:
// Info — для обычной информации
// Error — для ошибок
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// logger — конкретная реализация интерфейса Logger.
// Внутри два объекта log.Logger: один для stdout, другой для stderr.
type logger struct {
	info  *log.Logger
	error *log.Logger
}

// New создаёт новый логгер.
// Info пишет в стандартный вывод (os.Stdout), а Error — в стандартный поток ошибок (os.Stderr).
// log.LstdFlags добавляет дату и время к каждому сообщению.
func New() Logger {
	return &logger{
		info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// Info выводит информационное сообщение.
// Поддерживает форматирование (fmt.Sprintf-подобно).
func (l *logger) Info(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Error выводит сообщение об ошибке.
// Поддерживает форматирование (fmt.Sprintf-подобно).
func (l *logger) Error(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}
