package usecases

import (
	"context"
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

func (uc *CreateTeamUC) Invoke(ctx context.Context, input team.Entity) (*team.Entity, error) {
	err := uc.validateTeamMembers(ctx, input)
	if err != nil {
		return nil, err
	}

	err = uc.teamRepository.Save(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error while inserting team in database: %w", err)
	}
	// Return a pointer to the input entity
	return &input, nil
}

func (uc *CreateTeamUC) validateTeamMembers(ctx context.Context, input team.Entity) error {
	if len(input.Members) > 0 {
		playerIDs := make([]string, len(input.Members))
		for i, member := range input.Members {
			playerIDs[i] = member.ID
		}

		players, err := uc.playerRepository.Find(ctx, player.DomainQuery{
			Ids: playerIDs,
		})
		if err != nil {
			return fmt.Errorf("error while finding players: %w", err)
		}

		if len(players) != len(input.Members) {
			return fmt.Errorf("some of the team member does not exist")
		}
	}
	return nil
}
