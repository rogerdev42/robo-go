package models

import (
	"time"
)

type Category struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"-" db:"user_id"` // Не выводим в JSON
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	NotesCount int       `json:"notes_count" db:"notes_count"` // Вычисляемое поле
}

type CategoryCreate struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type CategoryUpdate struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}