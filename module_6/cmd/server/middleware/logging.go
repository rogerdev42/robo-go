package middleware

import (
	"module_6/internal/logger"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LoggingMiddleware creates middleware for HTTP request logging
func LoggingMiddleware(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		reqID := ""
		if id := c.Locals("requestid"); id != nil {
			reqID = id.(string)
		}

		if reqID != "" {
			c.Locals("logger", log.With(logger.String("request_id", reqID)))
		} else {
			c.Locals("logger", log)
		}

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		fields := []logger.Field{
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.Int("status", status),
			logger.Int64("duration_ms", duration.Milliseconds()),
			logger.String("ip", c.IP()),
		}

		if reqID != "" {
			fields = append(fields, logger.String("request_id", reqID))
		}

		if err != nil {
			fields = append(fields, logger.Error(err))
		}

		log.Info("request completed", fields...)

		return err
	}
}