package models

import "time"

// User represents a system user
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserCreate represents user registration request
type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLogin represents user login request
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
