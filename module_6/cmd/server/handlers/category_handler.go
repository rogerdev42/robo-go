package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// CategoryHandler handles category operations
type CategoryHandler struct {
	*BaseHandler
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(categoryService *services.CategoryService, log logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		BaseHandler:     NewBaseHandler(log),
		categoryService: categoryService,
	}
}

// GetAll gets all user categories
func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	categories, err := h.categoryService.GetAll(c.Context(), userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"categories": categories,
	})
}

// GetByID gets category by ID
func (h *CategoryHandler) GetByID(c fiber.Ctx) error {
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	category, err := h.categoryService.GetByID(c.Context(), categoryID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, category)
}

// Create creates new category
func (h *CategoryHandler) Create(c fiber.Ctx) error {
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetCategoryCreate(c)

	category, err := h.categoryService.Create(c.Context(), userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Created(c, category)
}

// Update updates category
func (h *CategoryHandler) Update(c fiber.Ctx) error {
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetCategoryUpdate(c)

	category, err := h.categoryService.Update(c.Context(), categoryID, userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, category)
}

// Delete deletes category
func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	categoryID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	affectedNotes, err := h.categoryService.Delete(c.Context(), categoryID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"message":        "Category deleted",
		"affected_notes": affectedNotes,
	})
}