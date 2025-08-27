package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type logger struct {
	info  *log.Logger
	error *log.Logger
}

func New() Logger {
	return &logger{
		info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

func (l *logger) Info(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}
