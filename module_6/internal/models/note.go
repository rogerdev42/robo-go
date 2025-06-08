package models

import (
	"database/sql"
	"time"
)

// Note represents a user note
type Note struct {
	ID         int           `json:"id" db:"id"`
	UserID     int           `json:"-" db:"user_id"`
	CategoryID sql.NullInt64 `json:"-" db:"category_id"`
	Title      string        `json:"title" db:"title"`
	Content    string        `json:"content" db:"content"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at"`
	Category   *Category     `json:"category,omitempty"`
}

// NoteCreate represents note creation request
type NoteCreate struct {
	Title      string `json:"title" validate:"required,min=1,max=255"`
	Content    string `json:"content" validate:"required,min=1"`
	CategoryID *int   `json:"category_id" validate:"omitempty,min=1"`
}

// NoteUpdate represents note update request
type NoteUpdate struct {
	Title      *string `json:"title" validate:"omitempty,min=1,max=255"`
	Content    *string `json:"content" validate:"omitempty,min=1"`
	CategoryID *int    `json:"category_id" validate:"omitempty"`
}

// NoteFilter represents note filtering options
type NoteFilter struct {
	CategoryID *int
	Search     string
	Limit      int
	Offset     int
}
