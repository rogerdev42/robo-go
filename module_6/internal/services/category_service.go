package services

import (
	"context"
	"errors"
	"module_6/internal/database/repository"
	"module_6/internal/logger"
	"module_6/internal/models"
)

// CategoryService сервис для работы с категориями
type CategoryService struct {
	categoryRepo repository.CategoryRepository
	logger       logger.Logger
}

// NewCategoryService создает новый экземпляр CategoryService
func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	log logger.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       log,
	}
}

// Create создает новую категорию
func (s *CategoryService) Create(ctx context.Context, userID int, req *models.CategoryCreate) (*models.Category, error) {
	s.logger.Info("Creating category",
		logger.Int("user_id", userID),
		logger.String("name", req.Name))

	// Проверяем, существует ли категория с таким именем у пользователя
	existingCategory, err := s.categoryRepo.GetByUserIDAndName(ctx, userID, req.Name)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		s.logger.Error("Failed to check existing category", logger.Error(err))
		return nil, err
	}

	if existingCategory != nil {
		s.logger.Warn("Category already exists",
			logger.Int("user_id", userID),
			logger.String("name", req.Name))
		return nil, models.ErrAlreadyExists
	}

	// Создаем категорию
	category := &models.Category{
		UserID: userID,
		Name:   req.Name,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		s.logger.Error("Failed to create category", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Category created successfully",
		logger.Int("category_id", category.ID),
		logger.Int("user_id", userID))

	return category, nil
}

// GetAll получает все категории пользователя
func (s *CategoryService) GetAll(ctx context.Context, userID int) ([]*models.Category, error) {
	s.logger.Debug("Getting categories for user", logger.Int("user_id", userID))

	categories, err := s.categoryRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get categories", logger.Error(err))
		return nil, err
	}

	s.logger.Debug("Categories retrieved",
		logger.Int("user_id", userID),
		logger.Int("count", len(categories)))

	return categories, nil
}

// GetByID получает категорию по ID с проверкой владельца
func (s *CategoryService) GetByID(ctx context.Context, id, userID int) (*models.Category, error) {
	s.logger.Debug("Getting category",
		logger.Int("category_id", id),
		logger.Int("user_id", userID))

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.logger.Warn("Category not found", logger.Int("category_id", id))
		} else {
			s.logger.Error("Failed to get category", logger.Error(err))
		}
		return nil, err
	}

	// Проверяем владельца
	if category.UserID != userID {
		s.logger.Warn("Access forbidden to category",
			logger.Int("category_id", id),
			logger.Int("user_id", userID),
			logger.Int("owner_id", category.UserID))
		return nil, models.ErrForbidden
	}

	return category, nil
}

// Update обновляет категорию
func (s *CategoryService) Update(ctx context.Context, id, userID int, req *models.CategoryUpdate) (*models.Category, error) {
	s.logger.Info("Updating category",
		logger.Int("category_id", id),
		logger.Int("user_id", userID),
		logger.String("new_name", req.Name))

	// Получаем категорию с проверкой владельца
	category, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Если имя не изменилось - возвращаем как есть
	if category.Name == req.Name {
		return category, nil
	}

	// Проверяем, не занято ли новое имя
	existingCategory, err := s.categoryRepo.GetByUserIDAndName(ctx, userID, req.Name)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		s.logger.Error("Failed to check existing category", logger.Error(err))
		return nil, err
	}

	if existingCategory != nil && existingCategory.ID != id {
		s.logger.Warn("Category name already taken",
			logger.Int("user_id", userID),
			logger.String("name", req.Name))
		return nil, models.ErrAlreadyExists
	}

	// Обновляем категорию
	category.Name = req.Name
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		s.logger.Error("Failed to update category", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Category updated successfully",
		logger.Int("category_id", id),
		logger.Int("user_id", userID))

	return category, nil
}

// Delete удаляет категорию
func (s *CategoryService) Delete(ctx context.Context, id, userID int) (int64, error) {
	s.logger.Info("Deleting category",
		logger.Int("category_id", id),
		logger.Int("user_id", userID))

	// Проверяем владельца
	category, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return 0, err
	}

	// Обновляем заметки - убираем у них категорию
	affectedNotes, err := s.categoryRepo.UpdateNotesCategory(ctx, category.ID, userID)
	if err != nil {
		s.logger.Error("Failed to update notes", logger.Error(err))
		return 0, err
	}

	// Удаляем категорию
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete category", logger.Error(err))
		return 0, err
	}

	s.logger.Info("Category deleted successfully",
		logger.Int("category_id", id),
		logger.Int("user_id", userID),
		logger.Int64("affected_notes", affectedNotes))

	return affectedNotes, nil
}
