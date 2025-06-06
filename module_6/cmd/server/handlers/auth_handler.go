package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// AuthHandler handles authentication
type AuthHandler struct {
	*BaseHandler
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *services.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(log),
		authService: authService,
	}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	req := middleware.GetUserCreate(c)

	user, token, err := h.authService.SignUp(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Created(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}

// SignIn handles user login
func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	req := middleware.GetUserLogin(c)

	user, token, err := h.authService.SignIn(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}