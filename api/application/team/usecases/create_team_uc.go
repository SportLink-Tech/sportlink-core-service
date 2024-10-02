package usecases

import (
	"fmt"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
)

type CreateTeamUC struct {
	playerRepository player.Repository
	teamRepository   team.Repository
}

func NewCreateTeamUC(playerRepository player.Repository, teamRepository team.Repository) *CreateTeamUC {
	return &CreateTeamUC{
		playerRepository: playerRepository,
		teamRepository:   teamRepository,
	}
}

func (uc *CreateTeamUC) Invoke(input team.Entity) (*team.Entity, error) {
	err := uc.teamRepository.Save(input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting team in database: %w", err)
	}
	return &input, nil
}
