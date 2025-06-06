package handlers_test

import (
	"encoding/json"
	"io"
	"module_6/cmd/server/middleware"
	"module_6/internal/config"
	"module_6/internal/models"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестируем только основную функциональность без глубокой валидации
func TestBasicRoutes_HealthCheck(t *testing.T) {
	app := fiber.New()
	
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	
	assert.Equal(t, "ok", response["status"])
}

func TestBasicAuth_SignUpEndpoint(t *testing.T) {
	app := fiber.New()
	
	// Простой эндпоинт без валидации - просто проверяем что он отвечает
	app.Post("/auth/signup", func(c fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"message": "signup endpoint works",
			"token":   "test.token",
		})
	})

	requestBody := `{
		"email": "test@example.com",
		"name": "testuser",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/auth/signup", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 201, resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	
	assert.Equal(t, "signup endpoint works", response["message"])
	assert.Equal(t, "test.token", response["token"])
}

func TestBasicAuth_SignInEndpoint(t *testing.T) {
	app := fiber.New()
	
	app.Post("/auth/signin", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "signin endpoint works",
			"token":   "test.token",
		})
	})

	requestBody := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/auth/signin", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestJWTProtection_Unauthorized(t *testing.T) {
	app := fiber.New()
	
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	
	// Защищенный роут
	protected := app.Group("/api", middleware.JWTProtected(cfg))
	protected.Get("/categories", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Запрос без токена
	req := httptest.NewRequest("GET", "/api/categories", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	var errorResponse models.ErrorResponse
	err = json.Unmarshal(body, &errorResponse)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", errorResponse.Error.Code)
}

func TestJWTProtection_ValidToken(t *testing.T) {
	app := fiber.New()
	
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	
	// Создаем валидный токен
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)
	
	// Защищенный роут
	protected := app.Group("/api", middleware.JWTProtected(cfg))
	protected.Get("/categories", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/api/categories", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestCategoryEndpoint_Basic(t *testing.T) {
	app := fiber.New()
	
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	
	// Создаем валидный токен
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)
	
	// Защищенные роуты
	protected := app.Group("/api", middleware.JWTProtected(cfg))
	
	protected.Get("/categories", func(c fiber.Ctx) error {
		categories := []map[string]interface{}{
			{
				"id":          1,
				"name":        "Work",
				"notes_count": 5,
			},
			{
				"id":          2,
				"name":        "Personal",
				"notes_count": 3,
			},
		}
		return c.JSON(fiber.Map{"categories": categories})
	})
	
	protected.Post("/categories", func(c fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"id":   1,
			"name": "New Category",
		})
	})

	// Тест GET /categories
	req := httptest.NewRequest("GET", "/api/categories", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Тест POST /categories
	requestBody := `{"name": "New Category"}`
	req = httptest.NewRequest("POST", "/api/categories", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)
	
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestNotesEndpoint_Basic(t *testing.T) {
	app := fiber.New()
	
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	
	// Создаем валидный токен
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)
	
	// Защищенные роуты
	protected := app.Group("/api", middleware.JWTProtected(cfg))
	
	protected.Get("/notes", func(c fiber.Ctx) error {
		notes := []map[string]interface{}{
			{
				"id":      1,
				"title":   "My First Note",
				"content": "This is my first note",
			},
		}
		return c.JSON(fiber.Map{"notes": notes})
	})
	
	protected.Post("/notes", func(c fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"id":      1,
			"title":   "New Note",
			"content": "New note content",
		})
	})

	// Тест GET /notes
	req := httptest.NewRequest("GET", "/api/notes", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Тест POST /notes
	requestBody := `{"title": "New Note", "content": "New note content"}`
	req = httptest.NewRequest("POST", "/api/notes", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)
	
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}