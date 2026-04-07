package account_test

import (
	"sportlink/api/domain/account"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Check(t *testing.T) {
	tests := []struct {
		name   string
		entity account.Entity
		then   func(t *testing.T, errors []error)
	}{
		{
			name:   "given valid account entity when checking then returns no errors",
			entity: account.Entity{Email: "test@example.com", Nickname: "testuser"},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name:   "given empty email when checking then returns email required error",
			entity: account.Entity{Email: "", Nickname: "testuser"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: email is required")
			},
		},
		{
			name:   "given invalid email format when checking then returns invalid email format error",
			entity: account.Entity{Email: "invalid-email", Nickname: "testuser"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
			},
		},
		{
			name:   "given email too long when checking then returns email too long error",
			entity: account.Entity{Email: strings.Repeat("a", 250) + "@example.com", Nickname: "testuser"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: email is too long (max 254 characters)")
			},
		},
		{
			name:   "given empty nickname when checking then returns no error",
			entity: account.Entity{Email: "test@example.com", Nickname: ""},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 0)
			},
		},
		{
			name:   "given nickname too short when checking then returns nickname length error",
			entity: account.Entity{Email: "test@example.com", Nickname: "ab"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname must be at least 3 characters long")
			},
		},
		{
			name:   "given nickname too long when checking then returns nickname too long error",
			entity: account.Entity{Email: "test@example.com", Nickname: strings.Repeat("a", 51)},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname is too long (max 50 characters)")
			},
		},
		{
			name:   "given nickname with invalid characters when checking then returns invalid characters error",
			entity: account.Entity{Email: "test@example.com", Nickname: "test@user!"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname can only contain letters, numbers, spaces, hyphens, and underscores")
			},
		},
		{
			name:   "given multiple validation errors when checking then returns all errors",
			entity: account.Entity{Email: "invalid-email", Nickname: "ab"},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 2)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
				assert.Contains(t, errors[1].Error(), "nickname: nickname must be at least 3 characters long")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := account.NewValidator()
			errors := validator.Check(tt.entity)
			tt.then(t, errors)
		})
	}
}
