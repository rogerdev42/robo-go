package handlers

import (
	"module_6/cmd/server/middleware"
	"module_6/internal/logger"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v3"
)

// NoteHandler обработчик для работы с заметками
type NoteHandler struct {
	*BaseHandler
	noteService *services.NoteService
}

// NewNoteHandler создает новый экземпляр NoteHandler
func NewNoteHandler(noteService *services.NoteService, log logger.Logger) *NoteHandler {
	return &NoteHandler{
		BaseHandler: NewBaseHandler(log),
		noteService: noteService,
	}
}

// GetAll получает все заметки пользователя
func (h *NoteHandler) GetAll(c fiber.Ctx) error {
	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Парсим фильтры из query параметров (с автоматической валидацией)
	filter := h.paramsHelper.GetNoteFilter(c)

	// Получаем заметки
	notes, err := h.noteService.GetAll(c.Context(), userID, filter)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем результат
	return h.responseHelper.Success(c, fiber.Map{
		"notes": notes,
	})
}

// GetByID получает заметку по ID
func (h *NoteHandler) GetByID(c fiber.Ctx) error {
	// Получаем ID заметки с автоматической валидацией
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем заметку
	note, err := h.noteService.GetByID(c.Context(), noteID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем заметку
	return h.responseHelper.Success(c, note)
}

// Create создает новую заметку
func (h *NoteHandler) Create(c fiber.Ctx) error {
	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем уже валидированные данные из middleware
	req := middleware.GetNoteCreate(c)

	// Создаем заметку
	note, err := h.noteService.Create(c.Context(), userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем созданную заметку
	return h.responseHelper.Created(c, note)
}

// Update обновляет заметку
func (h *NoteHandler) Update(c fiber.Ctx) error {
	// Получаем ID заметки с автоматической валидацией
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Получаем уже валидированные данные из middleware
	req := middleware.GetNoteUpdate(c)

	// Обновляем заметку
	note, err := h.noteService.Update(c.Context(), noteID, userID, &req)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем обновленную заметку
	return h.responseHelper.Success(c, note)
}

// Delete удаляет заметку
func (h *NoteHandler) Delete(c fiber.Ctx) error {
	// Получаем ID заметки с автоматической валидацией
	noteID, err := h.paramsHelper.GetIDParam(c, "id")
	if err != nil {
		return err
	}

	// Получаем user_id с автоматической обработкой ошибок
	userID, err := h.paramsHelper.GetUserID(c)
	if err != nil {
		return err
	}

	// Удаляем заметку
	err = h.noteService.Delete(c.Context(), noteID, userID)
	if err != nil {
		return h.responseHelper.HandleServiceError(c, err)
	}

	// Возвращаем успешный ответ
	return h.responseHelper.NoContent(c)
}
