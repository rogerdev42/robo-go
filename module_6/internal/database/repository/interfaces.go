package repository

import (
	"context"
	"module_6/internal/models"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByName(ctx context.Context, name string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

// CategoryRepository интерфейс для работы с категориями
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id int) (*models.Category, error)
	GetByUserID(ctx context.Context, userID int) ([]*models.Category, error)
	GetByUserIDAndName(ctx context.Context, userID int, name string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
	UpdateNotesCategory(ctx context.Context, categoryID int, userID int) (int64, error)
}

// NoteRepository интерфейс для работы с заметками
type NoteRepository interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id int) (*models.Note, error)
	GetByUserID(ctx context.Context, userID int, filter models.NoteFilter) ([]*models.Note, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context, userID int, filter models.NoteFilter) (int, error)
}
