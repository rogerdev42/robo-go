package middleware_test

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/config"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestToken(secret string, userID int, expireHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestJWTProtected_ValidToken(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	// Создаем валидный токен
	token, err := createTestToken(cfg.JWTSecret, 123, 24)
	require.NoError(t, err)

	// Применяем middleware
	app.Use(middleware.JWTProtected(cfg))

	// Тестовый роут
	app.Get("/test", func(c fiber.Ctx) error {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"user_id": userID})
	})

	// Создаем запрос с токеном
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestJWTProtected_MissingToken(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	app.Use(middleware.JWTProtected(cfg))
	app.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// НЕ добавляем Authorization header

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTProtected_InvalidTokenFormat(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	app.Use(middleware.JWTProtected(cfg))
	app.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTProtected_ExpiredToken(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	// Создаем истекший токен (истек час назад)
	claims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Истек час назад
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	require.NoError(t, err)

	app.Use(middleware.JWTProtected(cfg))
	app.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTProtected_InvalidSecret(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	// Создаем токен с другим секретом
	token, err := createTestToken("wrong-secret", 123, 24)
	require.NoError(t, err)

	app.Use(middleware.JWTProtected(cfg))
	app.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestGetUserID_Success(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c fiber.Ctx) error {
		// Симулируем что middleware уже установил user_id
		c.Locals("user_id", 456)

		userID, err := middleware.GetUserID(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetUserID_NotSet(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c fiber.Ctx) error {
		// НЕ устанавливаем user_id
		userID, err := middleware.GetUserID(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(fiber.Map{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}
