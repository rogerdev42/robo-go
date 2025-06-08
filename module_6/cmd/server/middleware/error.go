package middleware

import (
	"errors"
	"module_6/internal/logger"

	"github.com/gofiber/fiber/v3"
)

// ErrorHandler creates middleware for error handling
func ErrorHandler(log logger.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		log.Error("HTTP error",
			logger.Int("status", code),
			logger.String("path", c.Path()),
			logger.String("method", c.Method()),
			logger.Error(err),
		)

		return c.Status(code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    getErrorCode(code),
				"message": message,
			},
		})
	}
}

func getErrorCode(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "VALIDATION_ERROR"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	default:
		return "INTERNAL_ERROR"
	}
}
