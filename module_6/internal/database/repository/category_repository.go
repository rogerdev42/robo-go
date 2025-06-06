package repository

import (
	"context"
	"database/sql"
	"errors"
	"module_6/internal/models"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (user_id, name) 
		VALUES ($1, $2) 
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, category.UserID, category.Name).
		Scan(&category.ID, &category.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return models.ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	category := &models.Category{}
	query := `
		SELECT c.id, c.user_id, c.name, c.created_at,
			   COUNT(n.id) as notes_count
		FROM categories c
		LEFT JOIN notes n ON n.category_id = c.id
		WHERE c.id = $1
		GROUP BY c.id, c.user_id, c.name, c.created_at
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.NotesCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Category, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.created_at,
			   COUNT(n.id) as notes_count
		FROM categories c
		LEFT JOIN notes n ON n.category_id = c.id
		WHERE c.user_id = $1
		GROUP BY c.id, c.user_id, c.name, c.created_at
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.CreatedAt,
			&category.NotesCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) GetByUserIDAndName(ctx context.Context, userID int, name string) (*models.Category, error) {
	category := &models.Category{}
	query := `
		SELECT c.id, c.user_id, c.name, c.created_at,
			   COUNT(n.id) as notes_count
		FROM categories c
		LEFT JOIN notes n ON n.category_id = c.id
		WHERE c.user_id = $1 AND c.name = $2
		GROUP BY c.id, c.user_id, c.name, c.created_at
	`

	err := r.db.QueryRowContext(ctx, query, userID, name).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.NotesCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories 
		SET name = $1
		WHERE id = $2 AND user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, category.Name, category.ID, category.UserID)
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

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

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

func (r *categoryRepository) UpdateNotesCategory(ctx context.Context, categoryID int, userID int) (int64, error) {
	query := `
		UPDATE notes 
		SET category_id = NULL 
		WHERE category_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, categoryID, userID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}