package routes

import (
	"module_6/cmd/server/handlers"
	"module_6/cmd/server/middleware"
	"module_6/internal/config"
	"module_6/internal/logger"

	"github.com/gofiber/fiber/v3"
)

// Router handles route setup
type Router struct {
	app *fiber.App
	cfg *config.Config
	log logger.Logger
}

// New creates a new router
func New(app *fiber.App, cfg *config.Config, log logger.Logger) *Router {
	return &Router{
		app: app,
		cfg: cfg,
		log: log,
	}
}

// Setup configures all application routes
func (r *Router) Setup(
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	noteHandler *handlers.NoteHandler,
) {
	r.setupHealthCheck()
	r.setupAPIRoutes(authHandler, categoryHandler, noteHandler)
}

func (r *Router) setupHealthCheck() {
	r.app.Get("/health", handlers.HealthCheck)
}

func (r *Router) setupAPIRoutes(
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	noteHandler *handlers.NoteHandler,
) {
	api := r.app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/signin", authHandler.SignIn)

	protected := api.Group("/", middleware.JWTProtected(r.cfg))
	r.setupCategoryRoutes(protected, categoryHandler)
	r.setupNoteRoutes(protected, noteHandler)
}

func (r *Router) setupCategoryRoutes(protected fiber.Router, categoryHandler *handlers.CategoryHandler) {
	categories := protected.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetByID)
	categories.Post("/", categoryHandler.Create)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)
}

func (r *Router) setupNoteRoutes(protected fiber.Router, noteHandler *handlers.NoteHandler) {
	notes := protected.Group("/notes")
	notes.Get("/", noteHandler.GetAll)
	notes.Post("/", noteHandler.Create)
	notes.Get("/:id", noteHandler.GetByID)
	notes.Put("/:id", noteHandler.Update)
	notes.Delete("/:id", noteHandler.Delete)
}
