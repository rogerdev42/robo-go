package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// NoteHandler handles note operations
type NoteHandler struct {
	*BaseHandler
	noteService *services.NoteService
}

// NewNoteHandler creates a new NoteHandler instance
func NewNoteHandler(noteService *services.NoteService, log logger.Logger) *NoteHandler {
	return &NoteHandler{
		BaseHandler: NewBaseHandler(log),
		noteService: noteService,
	}
}

// GetAll gets all user notes
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

// GetByID gets note by ID
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

// Create creates new note
func (h *NoteHandler) Create(c fiber.Ctx) error {
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetNoteCreate(c)

	note, err := h.noteService.Create(c.Context(), userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Created(c, note)
}

// Update updates note
func (h *NoteHandler) Update(c fiber.Ctx) error {
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	req := middleware.GetNoteUpdate(c)

	note, err := h.noteService.Update(c.Context(), noteID, userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	return h.responseHelper.Success(c, note)
}

// Delete deletes note
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

	return h.responseHelper.NoContent(c)
}