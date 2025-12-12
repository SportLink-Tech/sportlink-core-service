package account

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
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

	// Validate email
	if err := v.validateEmail(entity.Email); err != nil {
		errs = append(errs, fmt.Errorf("email: %w", err))
	}

	// Validate nickname
	if err := v.validateNickname(entity.Nickname); err != nil {
		errs = append(errs, fmt.Errorf("nickname: %w", err))
	}

	// Validate password
	if err := v.validatePassword(entity.Password); err != nil {
		errs = append(errs, fmt.Errorf("password: %w", err))
	}

	return errs
}

func (v *validator) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// RFC 5322 simplified email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Additional checks
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

	// Allow alphanumeric, spaces, hyphens, and underscores
	nicknameRegex := regexp.MustCompile(`^[a-zA-Z0-9 _\-]+$`)
	if !nicknameRegex.MatchString(nickname) {
		return errors.New("nickname can only contain letters, numbers, spaces, hyphens, and underscores")
	}

	return nil
}

func (v *validator) validatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password is too long (max 128 characters)")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
