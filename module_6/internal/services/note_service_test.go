package services_test

import (
	"context"
	"database/sql"
	"errors"
	"module_6/internal/logger"
	"module_6/internal/models"
	"module_6/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockNoteRepository мок репозитория заметок
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Create(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	if args.Error(0) == nil {
		note.ID = 1
		note.CreatedAt = time.Now()
		note.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockNoteRepository) GetByID(ctx context.Context, id int) (*models.Note, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) GetByUserID(ctx context.Context, userID int, filter models.NoteFilter) ([]*models.Note, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	if args.Error(0) == nil {
		note.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNoteRepository) Count(ctx context.Context, userID int, filter models.NoteFilter) (int, error) {
	args := m.Called(ctx, userID, filter)
	return args.Int(0), args.Error(1)
}

func setupNoteService(t *testing.T) (*services.NoteService, *MockNoteRepository, *MockCategoryRepository) {
	mockNoteRepo := &MockNoteRepository{}
	mockCategoryRepo := &MockCategoryRepository{}
	
	log, err := logger.NewSlogLogger("debug", "text", "stdout")
	require.NoError(t, err)
	
	service := services.NewNoteService(mockNoteRepo, mockCategoryRepo, log)
	return service, mockNoteRepo, mockCategoryRepo
}

func TestNoteService_Create_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	req := &models.NoteCreate{
		Title:      "Test Note",
		Content:    "This is a test note",
		CategoryID: nil, // Без категории
	}
	userID := 1

	// Настраиваем моки
	mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Note")).
		Return(nil)

	// Выполняем тест
	note, err := service.Create(context.Background(), userID, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.Equal(t, req.Title, note.Title)
	assert.Equal(t, req.Content, note.Content)
	assert.Equal(t, userID, note.UserID)
	assert.False(t, note.CategoryID.Valid) // Нет категории
	assert.Nil(t, note.Category)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Create_WithCategory_Success(t *testing.T) {
	service, mockNoteRepo, mockCategoryRepo := setupNoteService(t)

	categoryID := 1
	req := &models.NoteCreate{
		Title:      "Test Note",
		Content:    "This is a test note",
		CategoryID: &categoryID,
	}
	userID := 1

	// Мок категории
	category := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "Test Category",
	}

	// Настраиваем моки
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).
		Return(category, nil)
	mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Note")).
		Return(nil)

	// Выполняем тест
	note, err := service.Create(context.Background(), userID, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.Equal(t, req.Title, note.Title)
	assert.Equal(t, req.Content, note.Content)
	assert.True(t, note.CategoryID.Valid)
	assert.Equal(t, int64(categoryID), note.CategoryID.Int64)
	assert.Equal(t, category, note.Category)

	mockNoteRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestNoteService_Create_CategoryNotFound(t *testing.T) {
	service, _, mockCategoryRepo := setupNoteService(t)

	categoryID := 999
	req := &models.NoteCreate{
		Title:      "Test Note",
		Content:    "This is a test note",
		CategoryID: &categoryID,
	}
	userID := 1

	// Категория не найдена
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).
		Return(nil, models.ErrNotFound)

	// Выполняем тест
	note, err := service.Create(context.Background(), userID, req)

	// Проверки
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrInvalidInput))
	assert.Nil(t, note)

	mockCategoryRepo.AssertExpectations(t)
}

func TestNoteService_Create_CategoryAccessDenied(t *testing.T) {
	service, _, mockCategoryRepo := setupNoteService(t)

	categoryID := 1
	req := &models.NoteCreate{
		Title:      "Test Note",
		Content:    "This is a test note",
		CategoryID: &categoryID,
	}
	userID := 1

	// Категория принадлежит другому пользователю
	category := &models.Category{
		ID:     categoryID,
		UserID: 2, // Другой пользователь
		Name:   "Other User Category",
	}

	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).
		Return(category, nil)

	// Выполняем тест
	note, err := service.Create(context.Background(), userID, req)

	// Проверки
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrForbidden))
	assert.Nil(t, note)

	mockCategoryRepo.AssertExpectations(t)
}

func TestNoteService_GetAll_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	userID := 1
	filter := models.NoteFilter{
		CategoryID: nil,
		Search:     "",
		Limit:      20,
		Offset:     0,
	}

	expectedNotes := []*models.Note{
		{
			ID:      1,
			UserID:  userID,
			Title:   "Note 1",
			Content: "Content 1",
		},
		{
			ID:      2,
			UserID:  userID,
			Title:   "Note 2",
			Content: "Content 2",
		},
	}

	mockNoteRepo.On("GetByUserID", mock.Anything, userID, filter).
		Return(expectedNotes, nil)

	// Выполняем тест
	notes, err := service.GetAll(context.Background(), userID, filter)

	// Проверки
	assert.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, expectedNotes, notes)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetByID_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 1
	userID := 1
	expectedNote := &models.Note{
		ID:      noteID,
		UserID:  userID,
		Title:   "Test Note",
		Content: "Test Content",
	}

	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(expectedNote, nil)

	// Выполняем тест
	note, err := service.GetByID(context.Background(), noteID, userID)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, note)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetByID_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 999
	userID := 1

	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(nil, models.ErrNotFound)

	// Выполняем тест
	note, err := service.GetByID(context.Background(), noteID, userID)

	// Проверки
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrNotFound))
	assert.Nil(t, note)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_GetByID_AccessDenied(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 1
	userID := 1
	otherUserNote := &models.Note{
		ID:     noteID,
		UserID: 2, // Другой пользователь
		Title:  "Other User Note",
	}

	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(otherUserNote, nil)

	// Выполняем тест
	note, err := service.GetByID(context.Background(), noteID, userID)

	// Проверки
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrNotFound)) // Возвращает NotFound для безопасности
	assert.Nil(t, note)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Update_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 1
	userID := 1
	newTitle := "Updated Title"
	newContent := "Updated Content"

	req := &models.NoteUpdate{
		Title:   &newTitle,
		Content: &newContent,
	}

	existingNote := &models.Note{
		ID:      noteID,
		UserID:  userID,
		Title:   "Old Title",
		Content: "Old Content",
	}

	// Настраиваем моки
	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(existingNote, nil)
	mockNoteRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Note")).
		Return(nil)

	// Выполняем тест
	note, err := service.Update(context.Background(), noteID, userID, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.Equal(t, newTitle, note.Title)
	assert.Equal(t, newContent, note.Content)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Update_WithCategory(t *testing.T) {
	service, mockNoteRepo, mockCategoryRepo := setupNoteService(t)

	noteID := 1
	userID := 1
	categoryID := 2

	req := &models.NoteUpdate{
		CategoryID: &categoryID,
	}

	existingNote := &models.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}

	category := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "New Category",
	}

	// Настраиваем моки
	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(existingNote, nil)
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).
		Return(category, nil)
	mockNoteRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Note")).
		Return(nil)

	// Выполняем тест
	note, err := service.Update(context.Background(), noteID, userID, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.True(t, note.CategoryID.Valid)
	assert.Equal(t, int64(categoryID), note.CategoryID.Int64)
	assert.Equal(t, category, note.Category)

	mockNoteRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestNoteService_Update_RemoveCategory(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 1
	userID := 1
	categoryID := 0 // 0 означает убрать категорию

	req := &models.NoteUpdate{
		CategoryID: &categoryID,
	}

	existingNote := &models.Note{
		ID:         noteID,
		UserID:     userID,
		Title:      "Test Note",
		CategoryID: sql.NullInt64{Valid: true, Int64: 5}, // Была категория
	}

	// Настраиваем моки
	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(existingNote, nil)
	mockNoteRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Note")).
		Return(nil)

	// Выполняем тест
	note, err := service.Update(context.Background(), noteID, userID, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.False(t, note.CategoryID.Valid) // Категория убрана
	assert.Nil(t, note.Category)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Delete_Success(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 1
	userID := 1

	existingNote := &models.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Note to Delete",
	}

	// Настраиваем моки
	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(existingNote, nil)
	mockNoteRepo.On("Delete", mock.Anything, noteID).
		Return(nil)

	// Выполняем тест
	err := service.Delete(context.Background(), noteID, userID)

	// Проверки
	assert.NoError(t, err)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteService_Delete_NotFound(t *testing.T) {
	service, mockNoteRepo, _ := setupNoteService(t)

	noteID := 999
	userID := 1

	mockNoteRepo.On("GetByID", mock.Anything, noteID).
		Return(nil, models.ErrNotFound)

	// Выполняем тест
	err := service.Delete(context.Background(), noteID, userID)

	// Проверки
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrNotFound))

	mockNoteRepo.AssertExpectations(t)
}