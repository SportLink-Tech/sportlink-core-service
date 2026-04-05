package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/account/usecases"
	"sportlink/api/domain/account"
	amocks "sportlink/mocks/api/domain/account"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreateAccountUC(t *testing.T) {
	ctx := context.Background()
	findErr := errors.New("database connection error")
	saveErr := errors.New("database error")

	tests := []struct {
		name  string
		input account.Entity
		given func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator)
		then  func(t *testing.T, result *account.Entity, err error)
	}{
		{
			name: "given valid account and email is available when invoke then saves and returns entity",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com" &&
						entity.Nickname == "jorge" &&
						entity.Password == "SportLink-Password1234"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "cabrerajjorge@gmail.com"
				})).Return([]account.Entity{}, nil)
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com" &&
						entity.Nickname == "jorge" &&
						entity.Password == "SportLink-Password1234"
				})).Return(nil)
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "cabrerajjorge@gmail.com", result.Email)
				assert.Equal(t, "jorge", result.Nickname)
				assert.Equal(t, "SportLink-Password1234", result.Password)
			},
		},
		{
			name: "given account already exists when invoke then returns error and does not save",
			input: account.Entity{
				Email:    "existing@example.com",
				Nickname: "existing",
				Password: "SportLink-Password1234",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				existingAccount := account.Entity{
					Email:    "existing@example.com",
					Nickname: "existing",
					Password: "$2a$10$existinghashedpasswordexample",
				}
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "existing@example.com"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "existing@example.com"
				})).Return([]account.Entity{existingAccount}, nil)
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "account already exist")
			},
		},
		{
			name: "given validation fails when invoke then returns validation error and does not touch repository",
			input: account.Entity{
				Email:    "invalid-email",
				Nickname: "ab",
				Password: "weak",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "invalid-email"
				})).Return([]error{
					errors.New("email: invalid email format"),
					errors.New("nickname: nickname must be at least 3 characters long"),
					errors.New("password: password must be at least 8 characters long"),
				})
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "validation failed")
				assert.Contains(t, err.Error(), "email: invalid email format")
				assert.Contains(t, err.Error(), "nickname: nickname must be at least 3 characters long")
				assert.Contains(t, err.Error(), "password: password must be at least 8 characters long")
			},
		},
		{
			name: "given find fails when invoke then returns wrapped error",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "cabrerajjorge@gmail.com"
				})).Return([]account.Entity{}, findErr)
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while checking if account exists")
				assert.ErrorIs(t, err, findErr)
			},
		},
		{
			name: "given save fails when invoke then returns wrapped error",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "cabrerajjorge@gmail.com"
				})).Return([]account.Entity{}, nil)
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return(saveErr)
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting account in database")
				assert.ErrorIs(t, err, saveErr)
			},
		},
		{
			name: "given valid account with special password characters when invoke then succeeds",
			input: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "P@ssw0rd!123",
			},
			given: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "test@example.com" &&
						entity.Password == "P@ssw0rd!123"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "test@example.com"
				})).Return([]account.Entity{}, nil)
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "test@example.com" &&
						entity.Nickname == "testuser" &&
						entity.Password == "P@ssw0rd!123"
				})).Return(nil)
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "test@example.com", result.Email)
				assert.Equal(t, "testuser", result.Nickname)
				assert.Equal(t, "P@ssw0rd!123", result.Password)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repository := &amocks.Repository{}
			validator := &amocks.Validator{}
			uc := usecases.NewCreateAccountUC(repository, validator)

			tt.given(t, repository, validator)

			result, err := uc.Invoke(ctx, tt.input)

			tt.then(t, result, err)
			repository.AssertExpectations(t)
			validator.AssertExpectations(t)
		})
	}
}
