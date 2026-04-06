package account

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Validator interface {
	Check(entity Entity) []error
}

type validator struct{}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Check(entity Entity) []error {
	var errs []error

	if err := v.validateEmail(entity.Email); err != nil {
		errs = append(errs, fmt.Errorf("email: %w", err))
	}

	if err := v.validateNickname(entity.Nickname); err != nil {
		errs = append(errs, fmt.Errorf("nickname: %w", err))
	}

	return errs
}

func (v *validator) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email is too long (max 254 characters)")
	}

	if strings.Count(email, "@") != 1 {
		return errors.New("email must contain exactly one @ symbol")
	}

	return nil
}

func (v *validator) validateNickname(nickname string) error {
	if nickname == "" {
		return errors.New("nickname is required")
	}

	if len(nickname) < 3 {
		return errors.New("nickname must be at least 3 characters long")
	}

	if len(nickname) > 50 {
		return errors.New("nickname is too long (max 50 characters)")
	}

	nicknameRegex := regexp.MustCompile(`^[a-zA-Z0-9 _\-]+$`)
	if !nicknameRegex.MatchString(nickname) {
		return errors.New("nickname can only contain letters, numbers, spaces, hyphens, and underscores")
	}

	return nil
}
