package middleware_test

import (
	"encoding/json"
	"io"
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/models"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupValidationTest(t *testing.T) (*fiber.App, logger.Logger) {
	app := fiber.New()
	log, err := logger.NewSlogLogger("debug", "text", "stdout")
	require.NoError(t, err)
	return app, log
}

func TestValidateUserCreate_Success(t *testing.T) {
	app, log := setupValidationTest(t)

	// Применяем middleware правильным способом
	app.Use("/test", middleware.ValidateUserCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		req := middleware.GetUserCreate(c)
		return c.JSON(fiber.Map{
			"email": req.Email,
			"name":  req.Name,
		})
	})

	validJSON := `{
		"email": "test@example.com",
		"name": "testuser",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(validJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	// Проверяем что данные корректно извлечены
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "test@example.com", result["email"])
	assert.Equal(t, "testuser", result["name"])
}

func TestValidateUserCreate_InvalidEmail(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateUserCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	invalidJSON := `{
		"email": "invalid-email",
		"name": "testuser",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)

	// Проверяем что ошибка содержит информацию об email
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(body, &errorResp)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", errorResp.Error.Code)
	assert.Contains(t, errorResp.Error.Details, "email")
}

func TestValidateUserCreate_MissingFields(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateUserCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	incompleteJSON := `{
		"email": "test@example.com"
	}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(incompleteJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)

	// Проверяем что ошибка содержит информацию о недостающих полях
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(body, &errorResp)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", errorResp.Error.Code)
	// Должны быть ошибки для name и password
	assert.Contains(t, errorResp.Error.Details, "name")
	assert.Contains(t, errorResp.Error.Details, "password")
}

func TestValidateCategoryCreate_Success(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateCategoryCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		req := middleware.GetCategoryCreate(c)
		return c.JSON(fiber.Map{"name": req.Name})
	})

	validJSON := `{"name": "Test Category"}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(validJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	// Проверяем что данные переданы
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Test Category", result["name"])
}

func TestValidateCategoryCreate_NameTooShort(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateCategoryCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	invalidJSON := `{"name": "ab"}`  // Меньше 3 символов

	req := httptest.NewRequest("POST", "/test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestValidateNoteCreate_Success(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateNoteCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		req := middleware.GetNoteCreate(c)
		return c.JSON(fiber.Map{
			"title":   req.Title,
			"content": req.Content,
		})
	})

	validJSON := `{
		"title": "Test Note",
		"content": "This is a test note content",
		"category_id": 1
	}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(validJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestValidateInvalidJSON(t *testing.T) {
	app, log := setupValidationTest(t)

	app.Use("/test", middleware.ValidateUserCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	invalidJSON := `{invalid json}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)

	// Проверяем что получили правильную ошибку
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(body, &errorResp)
	require.NoError(t, err)

	assert.Equal(t, "INVALID_JSON", errorResp.Error.Code)
}

func TestGetValidatedRequest_Success(t *testing.T) {
	app, log := setupValidationTest(t)

	var capturedName string

	app.Use("/test", middleware.ValidateCategoryCreate(log))
	app.Post("/test", func(c fiber.Ctx) error {
		req := middleware.GetCategoryCreate(c)
		capturedName = req.Name
		return c.JSON(fiber.Map{"name": req.Name})
	})

	validJSON := `{"name": "Test Category"}`

	req := httptest.NewRequest("POST", "/test", strings.NewReader(validJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test Category", capturedName)
}

func TestGetValidatedRequest_NotSet(t *testing.T) {
	app, _ := setupValidationTest(t)

	app.Post("/test", func(c fiber.Ctx) error {
		// Пытаемся получить данные без валидации
		req := middleware.GetCategoryCreate(c)
		
		// Должны получить zero value
		assert.Empty(t, req.Name)
		
		return c.JSON(fiber.Map{"message": "no validation"})
	})

	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name": "test"}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}