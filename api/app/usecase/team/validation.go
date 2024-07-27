package team

import (
	"github.com/go-playground/validator/v10"
	"sportlink/api/domain/common"
	"strconv"
)

// customCategoryValidation checks if the Category is among the accepted values.
func customCategoryValidation(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}

	// Get the field value as a pointer to an integer.
	categoryPtr, ok := fl.Field().Interface().(*int)
	if !ok || categoryPtr == nil {
		return false // Not an integer pointer or is nil, hence invalid.
	}

	// Create a Category instance from the pointer.
	category := common.Category(*categoryPtr)

	// Use the String() method to validate if the value is part of the defined constants.
	// Convert back to int to ensure the value is not just an out of range integer.
	if catNum, err := strconv.Atoi(category.String()); err == nil {
		return common.Category(catNum) == category
	}
	return false
}
