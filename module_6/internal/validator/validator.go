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

// Get returns singleton validator instance with custom configurations
func Get() *validator.Validate {
	once.Do(func() {
		instance = validator.New()

		// Use JSON tags as field names in validation errors
		instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom username validation
		instance.RegisterValidation("username", func(fl validator.FieldLevel) bool {
			username := fl.Field().String()
			if len(username) < 3 || len(username) > 50 {
				return false
			}
			// Allow only letters, numbers and underscores
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

// Validate validates a struct using the configured validator
func Validate(s interface{}) error {
	return Get().Struct(s)
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return Get().Var(field, tag)
}
