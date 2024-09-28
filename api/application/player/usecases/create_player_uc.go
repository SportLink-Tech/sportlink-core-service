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

<<<<<<< HEAD
func (uc *CreatePlayerUC) Invoke(input player.Entity) (*player.Entity, error) {
	err := uc.repository.Save(input)
=======
func (uc *CreatePlayerUc) Invoke(input player.Entity) (*player.Entity, error) {
<<<<<<< HEAD
	err := uc.repository.Save(input)
=======
	err := uc.repository.Insert(input)
>>>>>>> ede636af5f05fcc09b639d934c1122b83ee8747b
>>>>>>> 09c444c30f6289d58512d6730f3c615f44f2fcbe
	if err != nil {
		return nil, fmt.Errorf("error while inserting player in database: %w", err)
	}
	return &input, nil
}
