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

// ValidateUserCreate middleware for UserCreate validation
func ValidateUserCreate(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.UserCreate

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// ValidateUserLogin middleware for UserLogin validation
func ValidateUserLogin(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.UserLogin

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// ValidateCategoryCreate middleware for CategoryCreate validation
func ValidateCategoryCreate(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.CategoryCreate

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// ValidateCategoryUpdate middleware for CategoryUpdate validation
func ValidateCategoryUpdate(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.CategoryUpdate

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// ValidateNoteCreate middleware for NoteCreate validation
func ValidateNoteCreate(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.NoteCreate

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// ValidateNoteUpdate middleware for NoteUpdate validation
func ValidateNoteUpdate(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req models.NoteUpdate

		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		if err := validator.Validate(&req); err != nil {
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

// Getters for validated data
func GetUserCreate(c fiber.Ctx) models.UserCreate {
	if req, ok := c.Locals("validated_request").(models.UserCreate); ok {
		return req
	}
	return models.UserCreate{}
}

func GetUserLogin(c fiber.Ctx) models.UserLogin {
	if req, ok := c.Locals("validated_request").(models.UserLogin); ok {
		return req
	}
	return models.UserLogin{}
}

func GetCategoryCreate(c fiber.Ctx) models.CategoryCreate {
	if req, ok := c.Locals("validated_request").(models.CategoryCreate); ok {
		return req
	}
	return models.CategoryCreate{}
}

func GetCategoryUpdate(c fiber.Ctx) models.CategoryUpdate {
	if req, ok := c.Locals("validated_request").(models.CategoryUpdate); ok {
		return req
	}
	return models.CategoryUpdate{}
}

func GetNoteCreate(c fiber.Ctx) models.NoteCreate {
	if req, ok := c.Locals("validated_request").(models.NoteCreate); ok {
		return req
	}
	return models.NoteCreate{}
}

func GetNoteUpdate(c fiber.Ctx) models.NoteUpdate {
	if req, ok := c.Locals("validated_request").(models.NoteUpdate); ok {
		return req
	}
	return models.NoteUpdate{}
}

// formatValidationErrors formats validation errors
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