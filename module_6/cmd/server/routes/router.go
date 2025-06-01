package routes

import (
	"module_6/cmd/server/handlers"
	"module_6/cmd/server/middleware"
	"module_6/internal/config"
	"module_6/internal/logger"

	"github.com/gofiber/fiber/v3"
)

// Router отвечает за настройку маршрутов
type Router struct {
	app *fiber.App
	cfg *config.Config
	log logger.Logger
}

// New создает новый роутер
func New(app *fiber.App, cfg *config.Config, log logger.Logger) *Router {
	return &Router{
		app: app,
		cfg: cfg,
		log: log,
	}
}

// Setup настраивает все маршруты приложения
func (r *Router) Setup(
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	noteHandler *handlers.NoteHandler,
) {
	r.setupHealthCheck()
	r.setupAPIRoutes(authHandler, categoryHandler, noteHandler)
}

// setupHealthCheck настраивает health check endpoint
func (r *Router) setupHealthCheck() {
	r.app.Get("/health", handlers.HealthCheck)
}

// setupAPIRoutes настраивает API маршруты
func (r *Router) setupAPIRoutes(
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	noteHandler *handlers.NoteHandler,
) {
	// API группа
	api := r.app.Group("/api")

	// Auth роуты (публичные) с валидацией
	auth := api.Group("/auth")

	// ПРАВИЛЬНЫЙ СИНТАКСИС: отдельные вызовы
	auth.Use("/signup", middleware.ValidateUserCreate(r.log))
	auth.Post("/signup", authHandler.SignUp)

	auth.Use("/signin", middleware.ValidateUserLogin(r.log))
	auth.Post("/signin", authHandler.SignIn)

	// Защищенные роуты
	protected := api.Group("/", middleware.JWTProtected(r.cfg))
	r.setupCategoryRoutes(protected, categoryHandler)
	r.setupNoteRoutes(protected, noteHandler)
}

// setupCategoryRoutes настраивает маршруты категорий
func (r *Router) setupCategoryRoutes(protected fiber.Router, categoryHandler *handlers.CategoryHandler) {
	categories := protected.Group("/categories")

	// GET роуты БЕЗ валидации
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetByID)

	// POST с валидацией - ПРАВИЛЬНО
	categories.Post("/",
		middleware.ValidateCategoryCreate(r.log),
		categoryHandler.Create)

	// PUT с валидацией - ПРАВИЛЬНО
	categories.Put("/:id",
		middleware.ValidateCategoryUpdate(r.log),
		categoryHandler.Update)

	// DELETE БЕЗ валидации - ПРАВИЛЬНО
	categories.Delete("/:id", categoryHandler.Delete)
}

// setupNoteRoutes настраивает маршруты заметок
func (r *Router) setupNoteRoutes(protected fiber.Router, noteHandler *handlers.NoteHandler) {
	notes := protected.Group("/notes")

	// GET роуты БЕЗ валидации
	notes.Get("/", noteHandler.GetAll)
	notes.Get("/:id", noteHandler.GetByID)

	// POST с валидацией
	notes.Post("/",
		middleware.ValidateNoteCreate(r.log),
		noteHandler.Create)

	// PUT с валидацией
	notes.Put("/:id",
		middleware.ValidateNoteUpdate(r.log),
		noteHandler.Update)

	// DELETE БЕЗ валидации
	notes.Delete("/:id", noteHandler.Delete)
}
