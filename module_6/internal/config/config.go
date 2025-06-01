package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	Port string
	Env  string // development, production

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret      string
	JWTExpireHours int

	// Logging
	LogLevel  string
	LogFormat string
	LogOutput string // stdout, stderr или путь к файлу
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{
		// Server defaults
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "notes_db"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		LogOutput: getEnv("LOG_OUTPUT", "stdout"),
	}

	// Парсим JWT_EXPIRE_HOURS
	expireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRE_HOURS: %w", err)
	}
	cfg.JWTExpireHours = expireHours

	// Валидация обязательных полей
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate проверяет обязательные поля конфигурации
func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if c.Env != "development" && c.Env != "production" {
		return fmt.Errorf("ENV must be 'development' or 'production'")
	}

	c.LogLevel = strings.ToLower(c.LogLevel)
	if c.LogLevel != "debug" && c.LogLevel != "info" && c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	c.LogFormat = strings.ToLower(c.LogFormat)
	if c.LogFormat != "json" && c.LogFormat != "text" {
		return fmt.Errorf("LOG_FORMAT must be 'json' or 'text'")
	}

	return nil
}

// DatabaseDSN возвращает строку подключения к БД
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
