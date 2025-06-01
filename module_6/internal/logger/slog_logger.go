package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// SlogLogger реализация Logger через slog
type SlogLogger struct {
	logger *slog.Logger
	file   *os.File // для закрытия файла если используется
}

// NewSlogLogger создает новый логгер на основе slog
func NewSlogLogger(level, format, output string) (LoggerWithCloser, error) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		// Добавляем источник вызова в development режиме
		AddSource: level == "debug",
	}

	// Определяем куда писать логи
	var writer io.Writer
	var file *os.File

	switch strings.ToLower(output) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// Это путь к файлу
		// Создаем директорию если не существует
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Открываем файл для записи
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		file = f
		writer = f
	}

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	return &SlogLogger{
		logger: slog.New(handler),
		file:   file,
	}, nil
}

// Close закрывает файл логов если он используется
func (l *SlogLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug логирует сообщение уровня DEBUG
func (l *SlogLogger) Debug(msg string, fields ...Field) {
	l.log(slog.LevelDebug, msg, fields...)
}

// Info логирует сообщение уровня INFO
func (l *SlogLogger) Info(msg string, fields ...Field) {
	l.log(slog.LevelInfo, msg, fields...)
}

// Warn логирует сообщение уровня WARN
func (l *SlogLogger) Warn(msg string, fields ...Field) {
	l.log(slog.LevelWarn, msg, fields...)
}

// Error логирует сообщение уровня ERROR
func (l *SlogLogger) Error(msg string, fields ...Field) {
	l.log(slog.LevelError, msg, fields...)
}

// With создает новый логгер с дополнительными полями
func (l *SlogLogger) With(fields ...Field) Logger {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return &SlogLogger{
		logger: l.logger.With(args...),
		file:   l.file, // передаем ссылку на файл
	}
}

// log внутренний метод для логирования
func (l *SlogLogger) log(level slog.Level, msg string, fields ...Field) {
	attrs := make([]slog.Attr, len(fields))
	for i, f := range fields {
		attrs[i] = slog.Any(f.Key, f.Value)
	}
	l.logger.LogAttrs(context.Background(), level, msg, attrs...)
}
