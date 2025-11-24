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
	result, err := uc.repository.Find(player.DomainQuery{
		Id:       input.ID,
		Category: input.Category,
		Sport:    input.Sport,
	})

	// If find returns results without error, player already exists
	if err == nil && len(result) > 0 {
		return nil, fmt.Errorf("Player already exist: %+v", input)
	}

	// Save the new player
	err = uc.repository.Save(input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting player in database: %w", err)
	}
	
	return &input, nil
}
