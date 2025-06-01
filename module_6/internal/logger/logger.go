package logger

import "io"

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

// LoggerWithCloser расширенный интерфейс с возможностью закрытия
type LoggerWithCloser interface {
	Logger
	io.Closer
}

// Field представляет поле для структурированного логирования
type Field struct {
	Key   string
	Value interface{}
}

// Helper функции для создания полей

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
