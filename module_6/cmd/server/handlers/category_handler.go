package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// CategoryHandler обработчик для работы с категориями
type CategoryHandler struct {
	*BaseHandler
	categoryService *services.CategoryService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler
func NewCategoryHandler(categoryService *services.CategoryService, log logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		BaseHandler:     NewBaseHandler(log),
		categoryService: categoryService,
	}
}

// GetAll получает все категории пользователя
func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем категории
	categories, err := h.categoryService.GetAll(c.Context(), userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем результат
	return h.responseHelper.Success(c, fiber.Map{
		"categories": categories,
	})
}

// GetByID получает категорию по ID
func (h *CategoryHandler) GetByID(c fiber.Ctx) error {
	// Получаем ID категории с автоматической валидацией
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем категорию
	category, err := h.categoryService.GetByID(c.Context(), categoryID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем категорию
	return h.responseHelper.Success(c, category)
}

// Create создает новую категорию
func (h *CategoryHandler) Create(c fiber.Ctx) error {
	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем уже валидированные данные из middleware
	req := middleware.GetCategoryCreate(c)

	// Создаем категорию
	category, err := h.categoryService.Create(c.Context(), userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем созданную категорию
	return h.responseHelper.Created(c, category)
}

// Update обновляет категорию
func (h *CategoryHandler) Update(c fiber.Ctx) error {
	// Получаем ID категории с автоматической валидацией
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем уже валидированные данные из middleware
	req := middleware.GetCategoryUpdate(c)

	// Обновляем категорию
	category, err := h.categoryService.Update(c.Context(), categoryID, userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем обновленную категорию
	return h.responseHelper.Success(c, category)
}

// Delete удаляет категорию
func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	// Получаем ID категории с автоматической валидацией
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Удаляем категорию
	affectedNotes, err := h.categoryService.Delete(c.Context(), categoryID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем результат
	return h.responseHelper.Success(c, fiber.Map{
		"message":        "Категорію видалено",
		"affected_notes": affectedNotes,
	})
}
