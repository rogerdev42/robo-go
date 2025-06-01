package handlers

import (
	"module_6/internal/logger"
	"time"

	"github.com/gofiber/fiber/v3"
)

// getLogger получает логгер из контекста или возвращает дефолтный
func getLogger(c fiber.Ctx, defaultLogger logger.Logger) logger.Logger {
	if log, ok := c.Locals("logger").(logger.Logger); ok {
		return log
	}
	return defaultLogger
}

// HealthCheck обработчик для проверки состояния приложения
func HealthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now(),
	})
}

// BaseHandler базовый handler с общими зависимостями
type BaseHandler struct {
	logger         logger.Logger
	responseHelper *ResponseHelper
	paramsHelper   *ParamsHelper
}

// NewBaseHandler создает базовый handler
func NewBaseHandler(log logger.Logger) *BaseHandler {
	responseHelper := NewResponseHelper(log)
	paramsHelper := NewParamsHelper(responseHelper)

	return &BaseHandler{
		logger:         log,
		responseHelper: responseHelper,
		paramsHelper:   paramsHelper,
	}
}