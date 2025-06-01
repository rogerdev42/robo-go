package database

import (
	"context"
	"database/sql"
	"fmt"
	"module_6/internal/config"
	"module_6/internal/logger"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database обертка над sql.DB
type Database struct {
	*sql.DB
	logger logger.Logger
}

// New создает новое подключение к PostgreSQL
func New(cfg *config.Config, log logger.Logger) (*Database, error) {
	dsn := cfg.DatabaseDSN()

	log.Info("Connecting to database",
		logger.String("host", cfg.DBHost),
		logger.String("port", cfg.DBPort),
		logger.String("database", cfg.DBName),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Простые, но правильные настройки пула соединений
	db.SetMaxOpenConns(25)                 // Максимум открытых соединений
	db.SetMaxIdleConns(5)                  // Idle соединения
	db.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения
	db.SetConnMaxIdleTime(5 * time.Minute) // Время idle

	log.Info("Database connection pool configured",
		logger.Int("max_open_conns", 25),
		logger.Int("max_idle_conns", 5),
	)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected successfully")

	return &Database{
		DB:     db,
		logger: log,
	}, nil
}
