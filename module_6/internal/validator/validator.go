package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	instance *validator.Validate
)

// Get возвращает singleton экземпляр валидатора
func Get() *validator.Validate {
	once.Do(func() {
		instance = validator.New()

		// Используем json теги как имена полей в ошибках
		instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Кастомный валидатор для имени пользователя
		instance.RegisterValidation("username", func(fl validator.FieldLevel) bool {
			username := fl.Field().String()
			if len(username) < 3 || len(username) > 50 {
				return false
			}
			// Разрешаем только буквы, цифры и подчеркивания
			for _, char := range username {
				if !((char >= 'a' && char <= 'z') ||
					(char >= 'A' && char <= 'Z') ||
					(char >= '0' && char <= '9') ||
					char == '_') {
					return false
				}
			}
			return true
		})
	})

	return instance
}

// Validate выполняет валидацию структуры
func Validate(s interface{}) error {
	return Get().Struct(s)
}

// ValidateVar валидирует отдельную переменную
func ValidateVar(field interface{}, tag string) error {
	return Get().Var(field, tag)
}