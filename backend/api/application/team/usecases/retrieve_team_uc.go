package usecases

import (
	"context"
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"
)

type RetrieveTeamUC struct {
	teamRepository team.Repository
}

func NewRetrieveTeamUC(teamRepository team.Repository) *RetrieveTeamUC {
	return &RetrieveTeamUC{
		teamRepository: teamRepository,
	}
}

func (uc *RetrieveTeamUC) Invoke(ctx context.Context, id team.ID) (*team.Entity, error) {
	teams, err := uc.teamRepository.Find(ctx, team.DomainQuery{
		Name: id.Name,
		Sports: []common.Sport{
			id.Sport,
		},
	})
	if err != nil {
		return nil, err
	}
	return &teams[0], nil
}
