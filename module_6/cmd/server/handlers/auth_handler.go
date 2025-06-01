package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// AuthHandler обработчик для авторизации
type AuthHandler struct {
	*BaseHandler
	authService *services.AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService *services.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(log),
		authService: authService,
	}
}

// SignUp обработчик регистрации
func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	// Получаем уже валидированные данные из middleware
	req := middleware.GetUserCreate(c)

	// Создаем пользователя
	user, token, err := h.authService.SignUp(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем успешный ответ
	return h.responseHelper.Created(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}

// SignIn обработчик входа
func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	// Получаем уже валидированные данные из middleware
	req := middleware.GetUserLogin(c)

	// Выполняем вход
	user, token, err := h.authService.SignIn(c.Context(), &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем успешный ответ
	return h.responseHelper.Success(c, fiber.Map{
		"token": token,
		"user":  user.ToResponse(),
	})
}
