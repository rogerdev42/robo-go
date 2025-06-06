package app

import (
	"module_6/cmd/server/handlers"
	"module_6/internal/database/repository"
	"module_6/internal/services"
)

// Dependencies contains all application dependencies
type Dependencies struct {
	// Repositories
	UserRepo     repository.UserRepository
	CategoryRepo repository.CategoryRepository
	NoteRepo     repository.NoteRepository

	// Services
	AuthService     *services.AuthService
	CategoryService *services.CategoryService
	NoteService     *services.NoteService

	// Handlers
	AuthHandler     *handlers.AuthHandler
	CategoryHandler *handlers.CategoryHandler
	NoteHandler     *handlers.NoteHandler
}

// initDependencies initializes all application dependencies
func (a *App) initDependencies() *Dependencies {
	deps := &Dependencies{}

	deps.UserRepo = repository.NewUserRepository(a.database.DB)
	deps.CategoryRepo = repository.NewCategoryRepository(a.database.DB)
	deps.NoteRepo = repository.NewNoteRepository(a.database.DB)

	deps.AuthService = services.NewAuthService(deps.UserRepo, a.config, a.logger)
	deps.CategoryService = services.NewCategoryService(deps.CategoryRepo, a.logger)
	deps.NoteService = services.NewNoteService(deps.NoteRepo, deps.CategoryRepo, a.logger)

	deps.AuthHandler = handlers.NewAuthHandler(deps.AuthService, a.logger)
	deps.CategoryHandler = handlers.NewCategoryHandler(deps.CategoryService, a.logger)
	deps.NoteHandler = handlers.NewNoteHandler(deps.NoteService, a.logger)

	return deps
}