package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/player"
)

type CreatePlayerUC struct {
	repository player.Repository
}

func NewCreatePlayerUC(repository player.Repository) CreatePlayerUC {
	return CreatePlayerUC{
		repository: repository,
	}
}

func (uc *CreatePlayerUC) Invoke(ctx context.Context, input player.Entity) (*player.Entity, error) {
	// With ULID, each player gets a unique ID, so we don't need to check for duplicates by ID
	// If needed, we could check by Category and Sport, but typically players are unique entities

	// Save the new player
	err := uc.repository.Save(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting player in database: %w", err)
	}

	return &input, nil
}
