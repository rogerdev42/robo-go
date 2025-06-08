package models

import "time"

// Category represents a note category
type Category struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"-" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	NotesCount int       `json:"notes_count" db:"notes_count"`
}

// CategoryCreate represents category creation request
type CategoryCreate struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

// CategoryUpdate represents category update request
type CategoryUpdate struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
