package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"module_6/internal/models"
	"strings"
)

// noteRepository реализация репозитория для работы с заметками
type noteRepository struct {
	db *sql.DB
}

// NewNoteRepository создает новый экземпляр NoteRepository
func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

// Create создает новую заметку
func (r *noteRepository) Create(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (user_id, category_id, title, content) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`

	var categoryID sql.NullInt64
	if note.CategoryID.Valid {
		categoryID = note.CategoryID
	}

	err := r.db.QueryRowContext(ctx, query,
		note.UserID, categoryID, note.Title, note.Content).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID получает заметку по ID с информацией о категории
func (r *noteRepository) GetByID(ctx context.Context, id int) (*models.Note, error) {
	note := &models.Note{}
	query := `
		SELECT 
			n.id, n.user_id, n.category_id, n.title, n.content, 
			n.created_at, n.updated_at,
			c.id, c.name, c.created_at
		FROM notes n
		LEFT JOIN categories c ON n.category_id = c.id
		WHERE n.id = $1
	`

	var category sql.NullInt64
	var categoryName sql.NullString
	var categoryCreatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.UserID,
		&note.CategoryID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
		&category,
		&categoryName,
		&categoryCreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	// Если есть категория, добавляем ее в заметку
	if category.Valid {
		note.Category = &models.Category{
			ID:        int(category.Int64),
			Name:      categoryName.String,
			CreatedAt: categoryCreatedAt.Time,
		}
	}

	return note, nil
}

// GetByUserID получает все заметки пользователя с фильтрацией
func (r *noteRepository) GetByUserID(ctx context.Context, userID int, filter models.NoteFilter) ([]*models.Note, error) {
	// Базовый запрос
	query := `
		SELECT 
			n.id, n.user_id, n.category_id, n.title, n.content, 
			n.created_at, n.updated_at,
			c.id, c.name, c.created_at
		FROM notes n
		LEFT JOIN categories c ON n.category_id = c.id
		WHERE n.user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	// Добавляем фильтр по категории
	if filter.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND n.category_id = $%d", argCount)
		args = append(args, *filter.CategoryID)
	}

	// Добавляем поиск по заголовку и содержимому
	if filter.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (LOWER(n.title) LIKE LOWER($%d) OR LOWER(n.content) LIKE LOWER($%d))", argCount, argCount)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
	}

	// Сортировка
	query += " ORDER BY n.updated_at DESC"

	// Добавляем лимит и оффсет
	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		var categoryID sql.NullInt64
		var categoryName sql.NullString
		var categoryCreatedAt sql.NullTime

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.CategoryID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
			&categoryID,
			&categoryName,
			&categoryCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Если есть категория, добавляем ее в заметку
		if categoryID.Valid {
			note.Category = &models.Category{
				ID:        int(categoryID.Int64),
				Name:      categoryName.String,
				CreatedAt: categoryCreatedAt.Time,
			}
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// Update обновляет заметку
func (r *noteRepository) Update(ctx context.Context, note *models.Note) error {
	// Строим динамический запрос для обновления только переданных полей
	var setClauses []string
	var args []interface{}
	argCount := 0

	if note.Title != "" {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argCount))
		args = append(args, note.Title)
	}

	if note.Content != "" {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("content = $%d", argCount))
		args = append(args, note.Content)
	}

	// category_id может быть NULL
	argCount++
	setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argCount))
	if note.CategoryID.Valid {
		args = append(args, note.CategoryID.Int64)
	} else {
		args = append(args, nil)
	}

	// updated_at обновляется автоматически триггером

	// Добавляем WHERE условия
	argCount++
	args = append(args, note.ID)
	whereID := argCount

	argCount++
	args = append(args, note.UserID)
	whereUserID := argCount

	query := fmt.Sprintf(`
		UPDATE notes 
		SET %s
		WHERE id = $%d AND user_id = $%d
		RETURNING updated_at
	`, strings.Join(setClauses, ", "), whereID, whereUserID)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&note.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

// Delete удаляет заметку
func (r *noteRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM notes WHERE id = $1`

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

// Count возвращает количество заметок пользователя с учетом фильтров
func (r *noteRepository) Count(ctx context.Context, userID int, filter models.NoteFilter) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM notes n
		WHERE n.user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	// Добавляем фильтр по категории
	if filter.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND n.category_id = $%d", argCount)
		args = append(args, *filter.CategoryID)
	}

	// Добавляем поиск
	if filter.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (LOWER(n.title) LIKE LOWER($%d) OR LOWER(n.content) LIKE LOWER($%d))", argCount, argCount)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
