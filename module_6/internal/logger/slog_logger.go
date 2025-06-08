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

// SlogLogger implements Logger using slog
type SlogLogger struct {
	logger *slog.Logger
	file   *os.File
}

// NewSlogLogger creates a new slog-based logger
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
		Level:     logLevel,
		AddSource: level == "debug",
	}

	var writer io.Writer
	var file *os.File

	switch strings.ToLower(output) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

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

func (l *SlogLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *SlogLogger) Debug(msg string, fields ...Field) {
	l.log(slog.LevelDebug, msg, fields...)
}

func (l *SlogLogger) Info(msg string, fields ...Field) {
	l.log(slog.LevelInfo, msg, fields...)
}

func (l *SlogLogger) Warn(msg string, fields ...Field) {
	l.log(slog.LevelWarn, msg, fields...)
}

func (l *SlogLogger) Error(msg string, fields ...Field) {
	l.log(slog.LevelError, msg, fields...)
}

func (l *SlogLogger) With(fields ...Field) Logger {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return &SlogLogger{
		logger: l.logger.With(args...),
		file:   l.file,
	}
}

func (l *SlogLogger) log(level slog.Level, msg string, fields ...Field) {
	attrs := make([]slog.Attr, len(fields))
	for i, f := range fields {
		attrs[i] = slog.Any(f.Key, f.Value)
	}
	l.logger.LogAttrs(context.Background(), level, msg, attrs...)
}