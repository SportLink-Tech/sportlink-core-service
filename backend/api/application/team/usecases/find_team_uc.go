package usecases

import (
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

func (uc *FindTeamUC) Invoke(query team.DomainQuery) (*[]team.Entity, error) {
	result, err := uc.teamRepository.Find(query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
