package validator_test

import (
	"module_6/internal/models"
	"module_6/internal/validator"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreateValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    models.UserCreate
		wantErr bool
	}{
		{
			name: "valid user",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: models.UserCreate{
				Email:    "invalid-email",
				Name:     "testuser",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			user: models.UserCreate{
				Email:    "",
				Name:     "testuser",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short name",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     "ab",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "long name - 51 characters",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     strings.Repeat("a", 51), // Точно 51 символ
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "max valid name - 50 characters",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     strings.Repeat("a", 50), // Точно 50 символов
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "short password",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     "testuser",
				Password: "12345",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			user: models.UserCreate{
				Email:    "test@example.com",
				Name:     "testuser",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.user)
			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for name length: %d", len(tt.user.Name))
			} else {
				assert.NoError(t, err, "Expected no validation error for name length: %d", len(tt.user.Name))
			}
		})
	}
}

func TestCategoryCreateValidation(t *testing.T) {
	tests := []struct {
		name     string
		category models.CategoryCreate
		wantErr  bool
	}{
		{
			name: "valid category",
			category: models.CategoryCreate{
				Name: "Work",
			},
			wantErr: false,
		},
		{
			name: "short name",
			category: models.CategoryCreate{
				Name: "ab",
			},
			wantErr: true,
		},
		{
			name: "long name - 101 characters",
			category: models.CategoryCreate{
				Name: strings.Repeat("a", 101), // Точно 101 символ
			},
			wantErr: true,
		},
		{
			name: "max valid name - 100 characters",
			category: models.CategoryCreate{
				Name: strings.Repeat("a", 100), // Точно 100 символов
			},
			wantErr: false,
		},
		{
			name: "empty name",
			category: models.CategoryCreate{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.category)
			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for category name length: %d", len(tt.category.Name))
			} else {
				assert.NoError(t, err, "Expected no validation error for category name length: %d", len(tt.category.Name))
			}
		})
	}
}

func TestNoteCreateValidation(t *testing.T) {
	tests := []struct {
		name    string
		note    models.NoteCreate
		wantErr bool
	}{
		{
			name: "valid note without category",
			note: models.NoteCreate{
				Title:   "Test Note",
				Content: "Test content",
			},
			wantErr: false,
		},
		{
			name: "valid note with category",
			note: models.NoteCreate{
				Title:      "Test Note",
				Content:    "Test content",
				CategoryID: intPtr(1),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			note: models.NoteCreate{
				Title:   "",
				Content: "Test content",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			note: models.NoteCreate{
				Title:   "Test Note",
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "invalid category ID",
			note: models.NoteCreate{
				Title:      "Test Note",
				Content:    "Test content",
				CategoryID: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "title too long - 256 characters",
			note: models.NoteCreate{
				Title:   strings.Repeat("a", 256), // Точно 256 символов
				Content: "Test content",
			},
			wantErr: true,
		},
		{
			name: "max valid title - 255 characters",
			note: models.NoteCreate{
				Title:   strings.Repeat("a", 255), // Точно 255 символов
				Content: "Test content",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(&tt.note)
			if tt.wantErr {
				assert.Error(t, err, "Expected validation error")
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}