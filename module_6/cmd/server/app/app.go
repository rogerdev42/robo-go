package app

import (
	"context"
	"fmt"
	"module_6/cmd/server/middleware"
	"module_6/cmd/server/routes"
	"module_6/internal/config"
	"module_6/internal/database"
	"module_6/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// App представляет приложение
type App struct {
	config   *config.Config
	logger   logger.LoggerWithCloser
	database *database.Database
	fiber    *fiber.App
}

// New создает новое приложение
func New() (*App, error) {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Инициализируем логгер
	log, err := logger.NewSlogLogger(cfg.LogLevel, cfg.LogFormat, cfg.LogOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Подключаемся к базе данных
	db, err := database.New(cfg, log)
	if err != nil {
		log.Error("Failed to connect to database", logger.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	app := &App{
		config:   cfg,
		logger:   log,
		database: db,
	}

	// Создаем Fiber приложение
	app.fiber = app.createFiberApp()

	// Настраиваем middleware
	app.setupMiddleware()

	// Инициализируем зависимости и роуты
	deps := app.initDependencies()
	app.setupRoutes(deps)

	return app, nil
}

// Run запускает приложение
func (a *App) Run() error {
	a.logger.Info("Starting Notes API",
		logger.String("env", a.config.Env),
		logger.String("port", a.config.Port),
		logger.String("log_output", a.config.LogOutput),
	)

	// Канал для graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Запускаем сервер в горутине
	serverErrors := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%s", a.config.Port)
		a.logger.Info("Server starting", logger.String("address", addr))
		serverErrors <- a.fiber.Listen(addr)
	}()

	// Ждем либо ошибку сервера, либо сигнал завершения
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server failed: %w", err)

	case sig := <-shutdownChan:
		a.logger.Info("Shutdown signal received",
			logger.String("signal", sig.String()))

		// Graceful shutdown с таймаутом
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Останавливаем сервер
		if err := a.fiber.ShutdownWithContext(shutdownCtx); err != nil {
			a.logger.Error("Server shutdown error", logger.Error(err))
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		a.logger.Info("Server stopped gracefully")
		return nil
	}
}

// Close закрывает ресурсы приложения
func (a *App) Close() error {
	if a.database != nil {
		a.database.Close()
	}
	if a.logger != nil {
		a.logger.Close()
	}
	return nil
}

// createFiberApp создает экземпляр Fiber
func (a *App) createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:           "Notes API v1.0.0",
		StreamRequestBody: false,
		ServerHeader:      "Notes API",
		ErrorHandler:      middleware.ErrorHandler(a.logger),

		// Оптимизации производительности
		StrictRouting:     false,
		CaseSensitive:     false,
		UnescapePath:      false,
		BodyLimit:         4 * 1024 * 1024, // 4MB
		Concurrency:       256 * 1024,      // Максимум одновременных соединений
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		ReduceMemoryUsage: true, // Уменьшает использование памяти (~10% медленнее)
	})
}

// setupMiddleware настраивает middleware
func (a *App) setupMiddleware() {
	app := a.fiber

	// ТОЛЬКО БАЗОВЫЕ MIDDLEWARE
	app.Use(recover.New(recover.Config{
		EnableStackTrace: a.config.Env == "development",
	}))
	app.Use(requestid.New())
	app.Use(middleware.LoggingMiddleware(a.logger))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
}

// setupRoutes настраивает роуты приложения
func (a *App) setupRoutes(deps *Dependencies) {
	router := routes.New(a.fiber, a.config, a.logger)
	router.Setup(deps.AuthHandler, deps.CategoryHandler, deps.NoteHandler)
}
