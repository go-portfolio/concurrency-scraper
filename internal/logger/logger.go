package logger

// Logger — минимальный интерфейс логгера для приложения
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}
