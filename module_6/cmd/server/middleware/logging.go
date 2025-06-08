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

		// Получаем request ID более надежным способом
		reqID := ""
		if id := c.Locals("requestid"); id != nil {
			if idStr, ok := id.(string); ok {
				reqID = idStr
			}
		}

		// Создаем logger с request ID если он есть
		var contextLogger logger.Logger
		if reqID != "" {
			contextLogger = log.With(logger.String("request_id", reqID))
		} else {
			contextLogger = log
		}

		// Сохраняем logger в контексте
		c.Locals("logger", contextLogger)

		// Логируем начало запроса
		contextLogger.Info("request started",
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.String("ip", c.IP()),
			logger.String("user_agent", c.Get("User-Agent")),
		)

		// Выполняем следующий middleware/handler
		err := c.Next()

		// Вычисляем время выполнения
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Поля для логирования
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

		// Добавляем user_id если есть
		if userID, ok := c.Locals("user_id").(int); ok {
			fields = append(fields, logger.Int("user_id", userID))
		}

		// Логируем ошибку если есть
		if err != nil {
			fields = append(fields, logger.Error(err))
			contextLogger.Error("request failed", fields...)
		} else {
			// Логируем успешный запрос
			if status >= 400 {
				contextLogger.Warn("request completed with error", fields...)
			} else {
				contextLogger.Info("request completed", fields...)
			}
		}

		return err
	}
}
