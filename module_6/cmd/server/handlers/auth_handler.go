package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/models"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	*BaseHandler
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(log),
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	req := middleware.GetValidatedRequest[models.UserCreate](c)

	user, token, err := h.authService.SignUp(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Created(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}

func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	req := middleware.GetValidatedRequest[models.UserLogin](c)

	user, token, err := h.authService.SignIn(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}
