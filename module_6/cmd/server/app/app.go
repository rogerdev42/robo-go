package app

import (
	"context"
	"fmt"
	"module_6/cmd/server/middleware"
	"module_6/cmd/server/routes"
	"module_6/internal/config"
	"module_6/internal/database"
	"module_6/internal/logger"
	"module_6/internal/models"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// App represents the application
type App struct {
	config   *config.Config
	logger   logger.LoggerWithCloser
	database *database.Database
	fiber    *fiber.App
}

// New creates a new application
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log, err := logger.NewSlogLogger(cfg.LogLevel, cfg.LogFormat, cfg.LogOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

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

	app.fiber = app.createFiberApp()
	app.setupMiddleware()

	deps := app.initDependencies()
	app.setupRoutes(deps)

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	a.logger.Info("Starting Notes API",
		logger.String("env", a.config.Env),
		logger.String("port", a.config.Port),
		logger.String("log_output", a.config.LogOutput),
	)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	serverErrors := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%s", a.config.Port)
		a.logger.Info("Server starting", logger.String("address", addr))
		serverErrors <- a.fiber.Listen(addr)
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server failed: %w", err)

	case sig := <-shutdownChan:
		a.logger.Info("Shutdown signal received",
			logger.String("signal", sig.String()))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.fiber.ShutdownWithContext(shutdownCtx); err != nil {
			a.logger.Error("Server shutdown error", logger.Error(err))
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		a.logger.Info("Server stopped gracefully")
		return nil
	}
}

// Close closes application resources
func (a *App) Close() error {
	if a.database != nil {
		if err := a.database.Close(); err != nil {
			a.logger.Error("Failed to close database", logger.Error(err))
		}
	}
	if a.logger != nil {
		if err := a.logger.Close(); err != nil {
			fmt.Printf("Failed to close logger: %v\n", err)
		}
	}
	return nil
}

func (a *App) createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "Notes API v1.0.0",
		ErrorHandler: middleware.ErrorHandler(a.logger),
	})
}

func (a *App) setupMiddleware() {
	app := a.fiber

	// Сначала recover для обработки паник
	app.Use(recover.New(recover.Config{
		EnableStackTrace: a.config.Env == "development",
	}))

	// Затем request ID (важно что он идет ДО логирования)
	app.Use(requestid.New())

	// Логирование должно идти ПОСЛЕ request ID
	app.Use(middleware.LoggingMiddleware(a.logger))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Валидация для конкретных роутов
	app.Use("/api/auth/signup", middleware.ValidateJSON[models.UserCreate](a.logger))
	app.Use("/api/auth/signin", middleware.ValidateJSON[models.UserLogin](a.logger))

	app.Use("/api/categories", func(c fiber.Ctx) error {
		if c.Method() == "POST" {
			return middleware.ValidateJSON[models.CategoryCreate](a.logger)(c)
		}
		return c.Next()
	})

	app.Use("/api/categories/*", func(c fiber.Ctx) error {
		if c.Method() == "PUT" {
			return middleware.ValidateJSON[models.CategoryUpdate](a.logger)(c)
		}
		return c.Next()
	})

	app.Use("/api/notes", func(c fiber.Ctx) error {
		if c.Method() == "POST" {
			return middleware.ValidateJSON[models.NoteCreate](a.logger)(c)
		}
		return c.Next()
	})

	app.Use("/api/notes/*", func(c fiber.Ctx) error {
		if c.Method() == "PUT" {
			return middleware.ValidateJSON[models.NoteUpdate](a.logger)(c)
		}
		return c.Next()
	})
}

func (a *App) setupRoutes(deps *Dependencies) {
	router := routes.New(a.fiber, a.config, a.logger)
	router.Setup(deps.AuthHandler, deps.CategoryHandler, deps.NoteHandler)
}
