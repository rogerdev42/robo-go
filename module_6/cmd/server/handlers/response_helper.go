package handlers

import (
	"errors"
	"module_6/internal/logger"
	"module_6/internal/models"

	"github.com/gofiber/fiber/v3"
)

// ResponseHelper помощник для стандартизированных ответов
type ResponseHelper struct {
	logger logger.Logger
}

// NewResponseHelper создает новый helper
func NewResponseHelper(log logger.Logger) *ResponseHelper {
	return &ResponseHelper{logger: log}
}

// Success отправляет успешный ответ
func (h *ResponseHelper) Success(c fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

// Created отправляет ответ о создании ресурса
func (h *ResponseHelper) Created(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

// NoContent отправляет ответ без содержимого
func (h *ResponseHelper) NoContent(c fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// HandleServiceError обрабатывает ошибки из сервисов
func (h *ResponseHelper) HandleServiceError(c fiber.Ctx, err error) error {
	log := getLogger(c, h.logger)

	switch {
	case errors.Is(err, models.ErrNotFound):
		return h.errorResponse(c, fiber.StatusNotFound, "NOT_FOUND", "Resource not found")

	case errors.Is(err, models.ErrAlreadyExists):
		return h.errorResponse(c, fiber.StatusConflict, "ALREADY_EXISTS", "Resource already exists")

	case errors.Is(err, models.ErrForbidden):
		return h.errorResponse(c, fiber.StatusForbidden, "FORBIDDEN", "Access denied")

	case errors.Is(err, models.ErrInvalidInput):
		return h.errorResponse(c, fiber.StatusBadRequest, "INVALID_INPUT", "Invalid input data")

	case errors.Is(err, models.ErrInvalidCredentials):
		return h.errorResponse(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")

	default:
		log.Error("Unexpected service error", logger.Error(err))
		return h.errorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

// BadRequest отправляет ответ о некорректном запросе
func (h *ResponseHelper) BadRequest(c fiber.Ctx, message string) error {
	return h.errorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", message)
}

// Unauthorized отправляет ответ об отсутствии авторизации
func (h *ResponseHelper) Unauthorized(c fiber.Ctx, message string) error {
	return h.errorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message)
}

// errorResponse создает стандартный ответ с ошибкой
func (h *ResponseHelper) errorResponse(c fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}