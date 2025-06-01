package models

import "errors"

// Определяем типовые ошибки для приложения
var (
	// ErrNotFound возвращается когда запрошенный ресурс не найден
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists возвращается при попытке создать уже существующий ресурс
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrUnauthorized возвращается при отсутствии авторизации
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden возвращается при отсутствии прав доступа
	ErrForbidden = errors.New("access forbidden")

	// ErrInvalidInput возвращается при некорректных входных данных
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidCredentials возвращается при неверных учетных данных
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// ErrorResponse структура для ответа с ошибкой
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail детали ошибки
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
