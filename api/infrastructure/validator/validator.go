package validator

import (
	v10validator "github.com/go-playground/validator/v10"
	"sync"
)

var (
	once      sync.Once
	validator *v10validator.Validate
)

// GetInstance returns the singleton instance of the validator.
func GetInstance() *v10validator.Validate {
	once.Do(func() {
		validator = v10validator.New()
		validator.RegisterValidation("category", categoryValidation)
	})
	return validator
}
