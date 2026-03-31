package usecases

import (
	"context"
	"errors"
	"fmt"
	"sportlink/api/domain/account"
)

type CreateAccountUC struct {
	repository account.Repository
	validator  account.Validator
}

func NewCreateAccountUC(repository account.Repository, validator account.Validator) CreateAccountUC {
	return CreateAccountUC{
		repository: repository,
		validator:  validator,
	}
}

func (uc *CreateAccountUC) Invoke(ctx context.Context, input account.Entity) (*account.Entity, error) {
	// Validate entity
	if validationErrors := uc.validator.Check(input); len(validationErrors) > 0 {
		errMsg := "validation failed:"
		for _, err := range validationErrors {
			errMsg += fmt.Sprintf("\n  - %s", err.Error())
		}
		return nil, errors.New(errMsg)
	}

	// Ensure ID is generated using domain method
	if input.ID == "" {
		input = account.NewAccount(input.Email, input.Nickname, input.Password)
	}

	result, err := uc.repository.Find(ctx, account.DomainQuery{
		Emails: []string{input.Email},
	})

	if err != nil {
		return nil, fmt.Errorf("error while checking if account exists: %w", err)
	}

	if len(result) > 0 {
		return nil, fmt.Errorf("account already exist: %+v", input)
	}

	// Save the new account
	err = uc.repository.Save(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting account in database: %w", err)
	}

	return &input, nil
}
