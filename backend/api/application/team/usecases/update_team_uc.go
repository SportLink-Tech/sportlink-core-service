package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"
)

type UpdateTeamUC struct {
	teamRepository team.Repository
}

func NewUpdateTeamUC(teamRepository team.Repository) *UpdateTeamUC {
	return &UpdateTeamUC{teamRepository: teamRepository}
}

func (uc *UpdateTeamUC) Invoke(ctx context.Context, input team.PatchInput) (*team.Entity, error) {
	teams, err := uc.teamRepository.Find(ctx, team.DomainQuery{
		Name:   input.ID.Name,
		Sports: []common.Sport{input.ID.Sport},
	})
	if err != nil {
		return nil, fmt.Errorf("error finding team: %w", err)
	}
	if len(teams) == 0 {
		return nil, fmt.Errorf("team not found")
	}

	entity := teams[0]
	oldID := entity.ID

	if input.Name != nil {
		entity = entity.WithName(*input.Name)
	}

	if err = uc.teamRepository.Update(ctx, oldID, entity); err != nil {
		return nil, fmt.Errorf("error updating team: %w", err)
	}

	return &entity, nil
}
