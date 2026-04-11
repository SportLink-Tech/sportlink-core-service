package usecases

import (
	"context"
	"sportlink/api/domain/match"
)

type FindMatchesInput struct {
	AccountID string
	Statuses  []match.Status
}

type FindMatchesUC struct {
	matchRepository match.Repository
}

func NewFindMatchesUC(matchRepository match.Repository) *FindMatchesUC {
	return &FindMatchesUC{matchRepository: matchRepository}
}

func (uc *FindMatchesUC) Invoke(ctx context.Context, input FindMatchesInput) (*[]match.Entity, error) {
	entities, err := uc.matchRepository.Find(ctx, match.DomainQuery{
		AccountID: input.AccountID,
		Statuses:  input.Statuses,
	})
	if err != nil {
		return nil, err
	}
	return &entities, nil
}
