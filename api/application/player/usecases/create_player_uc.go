package usecases

import (
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

func (uc *CreatePlayerUC) Invoke(input player.Entity) (*player.Entity, error) {
	err := uc.repository.Save(input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting player in database: %w", err)
	}
	return &input, nil
}
