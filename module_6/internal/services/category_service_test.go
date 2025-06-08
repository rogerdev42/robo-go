package services_test

import (
	"context"
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

// MockCategoryRepository мок репозитория категорий
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	if args.Error(0) == nil {
		category.ID = 1
		category.CreatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Category, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByUserIDAndName(ctx context.Context, userID int, name string) (*models.Category, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) UpdateNotesCategory(ctx context.Context, categoryID int, userID int) (int64, error) {
	args := m.Called(ctx, categoryID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func setupCategoryService(t *testing.T) (*services.CategoryService, *MockCategoryRepository) {
	mockRepo := &MockCategoryRepository{}
	log, err := logger.NewSlogLogger("debug", "text", "stdout")
	require.NoError(t, err)
	
	service := services.NewCategoryService(mockRepo, log)
	return service, mockRepo
}

func TestCategoryService_Create_Success(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	req := &models.CategoryCreate{
		Name: "Test Category",
	}
	userID := 1

	// Категория с таким именем не существует
	mockRepo.On("GetByUserIDAndName", mock.Anything, userID, "Test Category").
		Return(nil, models.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Category")).
		Return(nil)

	category, err := service.Create(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Test Category", category.Name)
	assert.Equal(t, userID, category.UserID)
	assert.NotZero(t, category.ID)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create_AlreadyExists(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	req := &models.CategoryCreate{
		Name: "Existing Category",
	}
	userID := 1

	existingCategory := &models.Category{
		ID:     1,
		UserID: userID,
		Name:   "Existing Category",
	}

	mockRepo.On("GetByUserIDAndName", mock.Anything, userID, "Existing Category").
		Return(existingCategory, nil)

	category, err := service.Create(context.Background(), userID, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrAlreadyExists))
	assert.Nil(t, category)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetAll_Success(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	userID := 1
	expectedCategories := []*models.Category{
		{
			ID:         1,
			UserID:     userID,
			Name:       "Work",
			NotesCount: 5,
			CreatedAt:  time.Now(),
		},
		{
			ID:         2,
			UserID:     userID,
			Name:       "Personal",
			NotesCount: 3,
			CreatedAt:  time.Now(),
		},
	}

	mockRepo.On("GetByUserID", mock.Anything, userID).
		Return(expectedCategories, nil)

	categories, err := service.GetAll(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, expectedCategories, categories)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 1
	userID := 1
	expectedCategory := &models.Category{
		ID:         categoryID,
		UserID:     userID,
		Name:       "Test Category",
		NotesCount: 2,
		CreatedAt:  time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(expectedCategory, nil)

	category, err := service.GetByID(context.Background(), categoryID, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 999
	userID := 1

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(nil, models.ErrNotFound)

	category, err := service.GetByID(context.Background(), categoryID, userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrNotFound))
	assert.Nil(t, category)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID_AccessDenied(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 1
	userID := 1
	otherUserCategory := &models.Category{
		ID:     categoryID,
		UserID: 2, // Другой пользователь
		Name:   "Other User Category",
	}

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(otherUserCategory, nil)

	category, err := service.GetByID(context.Background(), categoryID, userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrForbidden))
	assert.Nil(t, category)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Success(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 1
	userID := 1
	req := &models.CategoryUpdate{
		Name: "Updated Category",
	}

	existingCategory := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "Old Category",
	}

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(existingCategory, nil)
	mockRepo.On("GetByUserIDAndName", mock.Anything, userID, "Updated Category").
		Return(nil, models.ErrNotFound)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Category")).
		Return(nil)

	category, err := service.Update(context.Background(), categoryID, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Updated Category", category.Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_NameTaken(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 1
	userID := 1
	req := &models.CategoryUpdate{
		Name: "Existing Name",
	}

	existingCategory := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "Old Category",
	}

	otherCategory := &models.Category{
		ID:     2,
		UserID: userID,
		Name:   "Existing Name",
	}

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(existingCategory, nil)
	mockRepo.On("GetByUserIDAndName", mock.Anything, userID, "Existing Name").
		Return(otherCategory, nil)

	category, err := service.Update(context.Background(), categoryID, userID, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrAlreadyExists))
	assert.Nil(t, category)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	service, mockRepo := setupCategoryService(t)

	categoryID := 1
	userID := 1
	expectedAffectedNotes := int64(5)

	existingCategory := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "Category to Delete",
	}

	mockRepo.On("GetByID", mock.Anything, categoryID).
		Return(existingCategory, nil)
	mockRepo.On("UpdateNotesCategory", mock.Anything, categoryID, userID).
		Return(expectedAffectedNotes, nil)
	mockRepo.On("Delete", mock.Anything, categoryID).
		Return(nil)

	affectedNotes, err := service.Delete(context.Background(), categoryID, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAffectedNotes, affectedNotes)

	mockRepo.AssertExpectations(t)
}