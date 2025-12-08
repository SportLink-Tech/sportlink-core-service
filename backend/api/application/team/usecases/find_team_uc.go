package usecases

import (
	"context"
	"sportlink/api/domain/team"
)

type FindTeamUC struct {
	teamRepository team.Repository
}

func NewFindTeamUC(teamRepository team.Repository) *FindTeamUC {
	return &FindTeamUC{
		teamRepository: teamRepository,
	}
}

func (uc *FindTeamUC) Invoke(ctx context.Context, query team.DomainQuery) (*[]team.Entity, error) {
	result, err := uc.teamRepository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
