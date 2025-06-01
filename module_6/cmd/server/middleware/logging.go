package middleware

import (
	"module_6/internal/logger"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LoggingMiddleware создает middleware для логирования HTTP запросов
func LoggingMiddleware(log logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Получаем request ID (может быть nil)
		reqID := ""
		if id := c.Locals("requestid"); id != nil {
			reqID = id.(string)
		}

		// Добавляем request ID в контекст для всех последующих логов
		if reqID != "" {
			c.Locals("logger", log.With(logger.String("request_id", reqID)))
		} else {
			c.Locals("logger", log)
		}

		// Обрабатываем запрос
		err := c.Next()

		// Логируем результат
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