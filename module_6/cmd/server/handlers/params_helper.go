package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// ParamsHelper helps with request parameter handling
type ParamsHelper struct {
	responseHelper *ResponseHelper
}

// NewParamsHelper creates new helper
func NewParamsHelper(responseHelper *ResponseHelper) *ParamsHelper {
	return &ParamsHelper{
		responseHelper: responseHelper,
	}
}

// GetIDParam extracts and validates ID from URL params
func (h *ParamsHelper) GetIDParam(c fiber.Ctx, paramName string) (int, error) {
	idStr := c.Params(paramName)
	if idStr == "" {
		return 0, h.responseHelper.BadRequest(c, "Missing ID parameter")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, h.responseHelper.BadRequest(c, "Invalid ID parameter")
	}

	return id, nil
}

// GetUserID extracts user ID from JWT token
func (h *ParamsHelper) GetUserID(c fiber.Ctx) (int, error) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return 0, h.responseHelper.Unauthorized(c, "User not authenticated")
	}
	return userID, nil
}

// GetNoteFilter parses note filters from query parameters
func (h *ParamsHelper) GetNoteFilter(c fiber.Ctx) models.NoteFilter {
	filter := models.NoteFilter{}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil && categoryID > 0 {
			filter.CategoryID = &categoryID
		}
	}

	filter.Search = c.Query("search")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		} else {
			filter.Limit = 20
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	return filter
}