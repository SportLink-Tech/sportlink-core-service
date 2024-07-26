package validator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	v10validator "github.com/go-playground/validator/v10"

	"sync"
)

var (
	once      sync.Once
	validator *v10validator.Validate
	uni       *ut.UniversalTranslator
)

// GetInstance returns the singleton instance of the validator.
func GetInstance() *v10validator.Validate {
	once.Do(func() {
		validator = v10validator.New()
		english := en.New()
		spanish := es.New()
		uni = ut.New(english, english, spanish)

		translator, _ := uni.GetTranslator("en")
		validator.RegisterValidation("category", customCategoryValidation)
		registerCustomTranslations(validator, translator)
	})
	return validator
}

func registerCustomTranslations(v *v10validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("category", trans, func(ut ut.Translator) error {
		return ut.Add("category", "'{0}' has a value of '{1}' which does not satisfy '{2}'.", true)
	}, func(ut ut.Translator, fe v10validator.FieldError) string {
		t, _ := ut.T("category", fe.Field(), fmt.Sprintf("%v", fe.Value()), fe.Tag())
		return t
	})
}

// customCategoryValidation checks if the Category is among the accepted values.
func customCategoryValidation(fl v10validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}

	category, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	allowedCategories := map[string]bool{"L1": true, "L2": true, "L3": true, "L4": true}
	if _, exists := allowedCategories[category]; exists {
		return true
	}
	return false
}
