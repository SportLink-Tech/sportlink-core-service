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
		translator := getTranslator()
		registerCustomTranslations(validator, translator)
	})
	return validator
}

func getTranslator() ut.Translator {
	english := en.New()
	spanish := es.New()
	uni = ut.New(english, english, spanish)
	translator, _ := uni.GetTranslator("en")
	return translator
}

func registerCustomTranslations(v *v10validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("category", trans, func(ut ut.Translator) error {
		return ut.Add("category", "'{0}' has a value of '{1}' which does not satisfy '{2}'.", true)
	}, func(ut ut.Translator, fe v10validator.FieldError) string {
		t, _ := ut.T("category", fe.Field(), fmt.Sprintf("%v", fe.Value()), fe.Tag())
		return t
	})
}
