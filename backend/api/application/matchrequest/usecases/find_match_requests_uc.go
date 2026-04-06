package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchrequest"
)

type FindMatchRequestsUC struct {
	matchRequestRepository matchrequest.Repository
}

func NewFindMatchRequestsUC(matchRequestRepository matchrequest.Repository) *FindMatchRequestsUC {
	return &FindMatchRequestsUC{
		matchRequestRepository: matchRequestRepository,
	}
}

func (uc *FindMatchRequestsUC) Invoke(ctx context.Context, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	entities, err := uc.matchRequestRepository.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while finding match requests: %w", err)
	}
	return entities, nil
}
