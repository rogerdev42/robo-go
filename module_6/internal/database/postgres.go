package database

import (
	"database/sql"
	"fmt"
	"module_6/internal/config"
	"module_6/internal/logger"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
	logger logger.Logger
}

func New(cfg *config.Config, log logger.Logger) (*Database, error) {
	dsn := cfg.DatabaseDSN()

	log.Info("Connecting to database")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected successfully")

	return &Database{
		DB:     db,
		logger: log,
	}, nil
}
