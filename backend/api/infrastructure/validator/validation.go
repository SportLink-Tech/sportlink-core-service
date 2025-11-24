package validator

import (
	v10validator "github.com/go-playground/validator/v10"
	"sportlink/api/domain/common"
)

func categoryValidation(fl v10validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}

	categoryVal, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}

	_, err := common.GetCategory(categoryVal)
	return err == nil
}
