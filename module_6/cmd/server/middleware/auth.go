package middleware

import (
	"module_6/internal/config"
	"module_6/internal/models"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected creates middleware for JWT token validation
func JWTProtected(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Missing authorization header",
				},
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Invalid authorization header format",
				},
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			if err == jwt.ErrTokenExpired {
				return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
					Error: models.ErrorDetail{
						Code:    "TOKEN_EXPIRED",
						Message: "Token is expired",
					},
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Invalid token",
				},
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Locals("user_id", int(userID))
				c.Locals("jwt", token)
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
					Error: models.ErrorDetail{
						Code:    "UNAUTHORIZED",
						Message: "Invalid token claims",
					},
				})
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Invalid token",
				},
			})
		}

		return c.Next()
	}
}

// GetUserID extracts user_id from context
func GetUserID(c fiber.Ctx) (int, error) {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return 0, models.ErrUnauthorized
	}
	return userID, nil
}