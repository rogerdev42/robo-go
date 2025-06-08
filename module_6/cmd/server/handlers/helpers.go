package handlers

import (
	"module_6/internal/logger"
	"time"

	"github.com/gofiber/fiber/v3"
)

// getLogger extracts logger from context or returns default
func getLogger(c fiber.Ctx, defaultLogger logger.Logger) logger.Logger {
	if log, ok := c.Locals("logger").(logger.Logger); ok {
		return log
	}
	return defaultLogger
}

// HealthCheck handles application health check
func HealthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now(),
	})
}

// BaseHandler contains common handler dependencies
type BaseHandler struct {
	logger         logger.Logger
	responseHelper *ResponseHelper
	paramsHelper   *ParamsHelper
}

// NewBaseHandler creates base handler
func NewBaseHandler(log logger.Logger) *BaseHandler {
	responseHelper := NewResponseHelper(log)
	paramsHelper := NewParamsHelper(responseHelper)

	return &BaseHandler{
		logger:         log,
		responseHelper: responseHelper,
		paramsHelper:   paramsHelper,
	}
}
