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
)

func TestCreateAccountUC_Invoke(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name  string
		input account.Entity
		on    func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator)
		then  func(t *testing.T, result *account.Entity, err error)
	}{
		{
			name: "given valid account when creating then returns created account",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
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
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "cabrerajjorge@gmail.com", result.Email)
				assert.Equal(t, "jorge", result.Nickname)
				assert.Equal(t, "SportLink-Password1234", result.Password)
			},
		},
		{
			name: "given account already exists when creating then returns error",
			input: account.Entity{
				Email:    "existing@example.com",
				Nickname: "existing",
				Password: "SportLink-Password1234",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
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
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "account already exist")
			},
		},
		{
			name: "given validation fails when creating then returns validation error",
			input: account.Entity{
				Email:    "invalid-email",
				Nickname: "ab",
				Password: "weak",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "invalid-email"
				})).Return([]error{
					errors.New("email: invalid email format"),
					errors.New("nickname: nickname must be at least 3 characters long"),
					errors.New("password: password must be at least 8 characters long"),
				})
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "validation failed")
				assert.Contains(t, err.Error(), "email: invalid email format")
				assert.Contains(t, err.Error(), "nickname: nickname must be at least 3 characters long")
				assert.Contains(t, err.Error(), "password: password must be at least 8 characters long")
			},
		},
		{
			name: "given repository find returns error when creating then returns error",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "cabrerajjorge@gmail.com"
				})).Return([]account.Entity{}, errors.New("database connection error"))
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while checking if account exists")
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
		{
			name: "given repository save fails when creating then returns error",
			input: account.Entity{
				Email:    "cabrerajjorge@gmail.com",
				Nickname: "jorge",
				Password: "SportLink-Password1234",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
				validator.On("Check", mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return([]error{})
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query account.DomainQuery) bool {
					return len(query.Emails) == 1 && query.Emails[0] == "cabrerajjorge@gmail.com"
				})).Return([]account.Entity{}, nil)
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity account.Entity) bool {
					return entity.Email == "cabrerajjorge@gmail.com"
				})).Return(errors.New("database error"))
			},
			then: func(t *testing.T, result *account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting account in database")
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "given valid account with special characters in password when creating then returns created account",
			input: account.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "P@ssw0rd!123",
			},
			on: func(t *testing.T, repository *amocks.Repository, validator *amocks.Validator) {
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
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "test@example.com", result.Email)
				assert.Equal(t, "testuser", result.Nickname)
				assert.Equal(t, "P@ssw0rd!123", result.Password)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			repository := amocks.NewRepository(t)
			validator := amocks.NewValidator(t)
			uc := usecases.NewCreateAccountUC(repository, validator)

			// given
			tt.on(t, repository, validator)

			// when
			result, err := uc.Invoke(ctx, tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
