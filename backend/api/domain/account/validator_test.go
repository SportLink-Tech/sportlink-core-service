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
			name: "given valid account entity when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given empty email when checking then returns email required error",
			entity: account.Entity{
				Email:    "",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: email is required")
			},
		},
		{
			name: "given invalid email format when checking then returns invalid email format error",
			entity: account.Entity{
				Email:    "invalid-email",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
			},
		},
		{
			name: "given email without domain when checking then returns invalid email format error",
			entity: account.Entity{
				Email:    "test@",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
			},
		},
		{
			name: "given email without @ symbol when checking then returns invalid email format error",
			entity: account.Entity{
				Email:    "testexample.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
			},
		},
		{
			name: "given email with multiple @ symbols when checking then returns invalid email format error",
			entity: account.Entity{
				Email:    "test@@example.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
			},
		},
		{
			name: "given email too long when checking then returns email too long error",
			entity: account.Entity{
				Email:    strings.Repeat("a", 250) + "@example.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "email: email is too long (max 254 characters)")
			},
		},
		{
			name: "given empty nickname when checking then returns nickname required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname is required")
			},
		},
		{
			name: "given nickname too short when checking then returns nickname length error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "ab",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname must be at least 3 characters long")
			},
		},
		{
			name: "given nickname too long when checking then returns nickname too long error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: strings.Repeat("a", 51),
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname is too long (max 50 characters)")
			},
		},
		{
			name: "given nickname with invalid characters when checking then returns invalid characters error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "test@user!",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "nickname: nickname can only contain letters, numbers, spaces, hyphens, and underscores")
			},
		},
		{
			name: "given valid nickname with spaces when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "test user",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given valid nickname with hyphens when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "test-user",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given valid nickname with underscores when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "test_user",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given empty password when checking then returns password required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password is required")
			},
		},
		{
			name: "given password too short when checking then returns password length error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "Short1!",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password must be at least 8 characters long")
			},
		},
		{
			name: "given password too long when checking then returns password too long error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: strings.Repeat("A", 129) + "1@a",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password is too long (max 128 characters)")
			},
		},
		{
			name: "given password without uppercase when checking then returns uppercase required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "lowercase123!",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password must contain at least one uppercase letter")
			},
		},
		{
			name: "given password without lowercase when checking then returns lowercase required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "UPPERCASE123!",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password must contain at least one lowercase letter")
			},
		},
		{
			name: "given password without number when checking then returns number required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "NoNumbers@!",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password must contain at least one number")
			},
		},
		{
			name: "given password without special character when checking then returns special character required error",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "NoSpecial123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 1)
				assert.Contains(t, errors[0].Error(), "password: password must contain at least one special character")
			},
		},
		{
			name: "given multiple validation errors when checking then returns all errors",
			entity: account.Entity{
				Email:    "invalid-email",
				Nickname: "ab",
				Password: "weak",
			},
			then: func(t *testing.T, errors []error) {
				assert.Len(t, errors, 3)
				assert.Contains(t, errors[0].Error(), "email: invalid email format")
				assert.Contains(t, errors[1].Error(), "nickname: nickname must be at least 3 characters long")
				assert.Contains(t, errors[2].Error(), "password: password must be at least 8 characters long")
			},
		},
		{
			name: "given password with punctuation character when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "ValidPass123!",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given password with symbol character when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "ValidPass123$",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given email with valid special characters when checking then returns no errors",
			entity: account.Entity{
				Email:    "test.user+tag@example.co.uk",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given minimum valid nickname when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "abc",
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given maximum valid nickname when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: strings.Repeat("a", 50),
				Password: "ValidP@ssw0rd123",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
		{
			name: "given minimum valid password when checking then returns no errors",
			entity: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "Valid1!a",
			},
			then: func(t *testing.T, errors []error) {
				assert.Empty(t, errors)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			validator := account.NewValidator()

			// when
			errors := validator.Check(tt.entity)

			// then
			tt.then(t, errors)
		})
	}
}
