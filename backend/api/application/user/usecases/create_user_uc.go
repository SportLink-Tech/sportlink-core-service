package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/user"
)

type CreateUserUC struct {
	repository user.Repository
}

func NewCreateUserUC(repository user.Repository) CreateUserUC {
	return CreateUserUC{
		repository: repository,
	}
}

func (uc *CreateUserUC) Invoke(ctx context.Context, input user.Entity) (*user.Entity, error) {
	// Save the new user
	err := uc.repository.Save(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting user in database: %w", err)
	}

	return &input, nil
}
