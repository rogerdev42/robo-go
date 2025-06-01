package repository

import (
	"context"
	"database/sql"
	"errors"
	"module_6/internal/models"
)

// UserRepository реализация репозитория для работы с пользователями
type userRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, name, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		// Проверяем на уникальность email или name
		if isUniqueViolation(err) {
			return models.ErrAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, password_hash, created_at 
		FROM users 
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, password_hash, created_at 
		FROM users 
		WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByName получает пользователя по имени
func (r *userRepository) GetByName(ctx context.Context, name string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, password_hash, created_at 
		FROM users 
		WHERE name = $1
	`

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update обновляет данные пользователя
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = $1, name = $2, password_hash = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Name, user.PasswordHash, user.ID)

	if err != nil {
		if isUniqueViolation(err) {
			return models.ErrAlreadyExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// Delete удаляет пользователя
func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}
