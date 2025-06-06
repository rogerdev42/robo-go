package handlers

import (
	"errors"
	"module_6/internal/logger"
	"module_6/internal/models"

	"github.com/gofiber/fiber/v3"
)

// ResponseHelper provides standardized responses
type ResponseHelper struct {
	logger logger.Logger
}

// NewResponseHelper creates new helper
func NewResponseHelper(log logger.Logger) *ResponseHelper {
	return &ResponseHelper{logger: log}
}

// Success sends successful response
func (h *ResponseHelper) Success(c fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

// Created sends resource creation response
func (h *ResponseHelper) Created(c fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

// NoContent sends no content response
func (h *ResponseHelper) NoContent(c fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// HandleServiceError handles service layer errors
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

// BadRequest sends bad request response
func (h *ResponseHelper) BadRequest(c fiber.Ctx, message string) error {
	return h.errorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", message)
}

// Unauthorized sends unauthorized response
func (h *ResponseHelper) Unauthorized(c fiber.Ctx, message string) error {
	return h.errorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message)
}

func (h *ResponseHelper) errorResponse(c fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}