package repository

import "strings"

// isUniqueViolation проверяет, является ли ошибка нарушением уникальности
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "duplicate key value violates unique constraint") ||
		strings.Contains(errMsg, "23505")
}
