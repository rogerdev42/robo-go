package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/models"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

type NoteHandler struct {
	*BaseHandler
	noteService *services.NoteService
}

func NewNoteHandler(noteService *services.NoteService, log logger.Logger) *NoteHandler {
	return &NoteHandler{
		BaseHandler: NewBaseHandler(log),
		noteService: noteService,
	}
}

func (h *NoteHandler) GetAll(c fiber.Ctx) error {
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	filter := h.paramsHelper.GetNoteFilter(c)

	notes, err := h.noteService.GetAll(c.Context(), userID, filter)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"notes": notes,
	})
}

func (h *NoteHandler) GetByID(c fiber.Ctx) error {
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	note, err := h.noteService.GetByID(c.Context(), noteID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, note)
}

func (h *NoteHandler) Create(c fiber.Ctx) error {
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetValidatedRequest[models.NoteCreate](c)

	note, err := h.noteService.Create(c.Context(), userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Created(c, note)
}

func (h *NoteHandler) Update(c fiber.Ctx) error {
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetValidatedRequest[models.NoteUpdate](c)

	note, err := h.noteService.Update(c.Context(), noteID, userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, note)
}

func (h *NoteHandler) Delete(c fiber.Ctx) error {
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	err = h.noteService.Delete(c.Context(), noteID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, fiber.Map{
		"message": "Category deleted",
	})
}
