package logger

// MockLogger
type MockLogger struct{}

func (l *MockLogger) Info(format string, args ...interface{})  {}
func (l *MockLogger) Error(format string, args ...interface{}) {}
func (l *MockLogger) Debug(format string, args ...interface{}) {}