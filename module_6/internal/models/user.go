package models

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"` // Не выводим в JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserCreate данные для создания пользователя
type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLogin данные для входа
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse ответ с данными пользователя
type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ToResponse преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}