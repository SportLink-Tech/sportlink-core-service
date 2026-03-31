package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/domain/matchrequest"
)

type CreateMatchRequestInput struct {
	MatchAnnouncementID string
	RequesterAccountID  string
}

type CreateMatchRequestUC struct {
	matchRequestRepository      matchrequest.Repository
	matchAnnouncementRepository matchannouncement.Repository
}

func NewCreateMatchRequestUC(
	matchRequestRepository matchrequest.Repository,
	matchAnnouncementRepository matchannouncement.Repository,
) *CreateMatchRequestUC {
	return &CreateMatchRequestUC{
		matchRequestRepository:      matchRequestRepository,
		matchAnnouncementRepository: matchAnnouncementRepository,
	}
}

func (uc *CreateMatchRequestUC) Invoke(ctx context.Context, input CreateMatchRequestInput) (*matchrequest.Entity, error) {
	// Fetch the match announcement to get the owner account ID
	page, err := uc.matchAnnouncementRepository.Find(ctx, matchannouncement.DomainQuery{
		IDs: []string{input.MatchAnnouncementID},
	})
	if err != nil {
		return nil, fmt.Errorf("error while finding match announcement: %w", err)
	}
	if len(page.Entities) == 0 {
		return nil, fmt.Errorf("match announcement '%s' not found", input.MatchAnnouncementID)
	}
	announcement := &page.Entities[0]

	// Prevent the owner from requesting their own announcement
	if announcement.OwnerAccountID == input.RequesterAccountID {
		return nil, fmt.Errorf("cannot send a match request to your own announcement")
	}

	entity := matchrequest.NewMatchRequest(
		input.MatchAnnouncementID,
		announcement.OwnerAccountID,
		input.RequesterAccountID,
	)

	if err := uc.matchRequestRepository.Save(ctx, entity); err != nil {
		return nil, fmt.Errorf("error while saving match request: %w", err)
	}

	return &entity, nil
}
