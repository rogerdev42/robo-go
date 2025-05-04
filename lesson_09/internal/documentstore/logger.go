package documentstore

import (
	"log/slog"
	"os"
)

var (
	l       *slog.Logger
	logFile *os.File
)

func init() {
	logFile, err := os.OpenFile("documentstore.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		l.Warn("failed to open log file. stdout is used", slog.String("error", err.Error()))
	} else {
		l = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
}

func CloseLogFile() {
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			l.Error("Failed to close log file", slog.String("error", err.Error()))
		}
	}
}
