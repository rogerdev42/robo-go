package middleware

import (
	"fmt"
	"module_6/internal/logger"
	"module_6/internal/models"
	"module_6/internal/validator"
	"strings"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// ValidateJSON validates JSON request body and stores validated data in context
func ValidateJSON[T any](log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req T

		if err := c.Bind().JSON(&req); err != nil {
			getLogger(c, log).Warn("Invalid JSON body", logger.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
			getLogger(c, log).Warn("Validation failed", logger.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "VALIDATION_ERROR",
					Message: "Validation failed",
					Details: formatValidationErrors(err),
				},
			})
		}

		c.Locals("validated_request", req)
		return c.Next()
	}
}

// GetValidatedRequest retrieves validated request data from context
func GetValidatedRequest[T any](c fiber.Ctx) T {
	if req, ok := c.Locals("validated_request").(T); ok {
		return req
	}
	var zero T
	return zero
}

// formatValidationErrors converts validator errors to user-friendly messages
func formatValidationErrors(err error) map[string]interface{} {
	errors := make(map[string]interface{})

	if validationErrors, ok := err.(validatorv10.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("Value is too short (minimum %s)", e.Param())
			case "max":
				errors[field] = fmt.Sprintf("Value is too long (maximum %s)", e.Param())
			case "username":
				errors[field] = "Username can only contain letters, numbers and underscores"
			default:
				errors[field] = "Invalid value"
			}
		}
	}

	return errors
}

func getLogger(c fiber.Ctx, defaultLogger logger.Logger) logger.Logger {
	if log, ok := c.Locals("logger").(logger.Logger); ok {
		return log
	}
	return defaultLogger
}
