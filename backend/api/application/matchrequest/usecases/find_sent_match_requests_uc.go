package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchrequest"
)

type FindSentMatchRequestsUC struct {
	matchRequestRepository matchrequest.Repository
}

func NewFindSentMatchRequestsUC(matchRequestRepository matchrequest.Repository) *FindSentMatchRequestsUC {
	return &FindSentMatchRequestsUC{
		matchRequestRepository: matchRequestRepository,
	}
}

func (uc *FindSentMatchRequestsUC) Invoke(ctx context.Context, requesterAccountID string, statuses []matchrequest.Status) ([]matchrequest.Entity, error) {
	entities, err := uc.matchRequestRepository.Find(ctx, matchrequest.DomainQuery{
		RequesterAccountIDs: []string{requesterAccountID},
		Statuses:            statuses,
	})
	if err != nil {
		return nil, fmt.Errorf("error while finding sent match requests: %w", err)
	}
	return entities, nil
}
