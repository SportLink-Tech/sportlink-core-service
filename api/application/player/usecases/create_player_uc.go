package usecases

import (
	"fmt"
	"sportlink/api/domain/player"
)

type CreatePlayerUc struct {
	repository player.Repository
}

func NewCreatePlayerUc(repository player.Repository) CreatePlayerUc {
	return CreatePlayerUc{
		repository: repository,
	}
}

func (uc *CreatePlayerUc) Invoke(input player.Entity) (*player.Entity, error) {
<<<<<<< HEAD
	err := uc.repository.Save(input)
=======
	err := uc.repository.Insert(input)
>>>>>>> ede636af5f05fcc09b639d934c1122b83ee8747b
	if err != nil {
		return nil, fmt.Errorf("error while inserting player in database: %w", err)
	}
	return &input, nil
}
