package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
)

type CreateMatchRequestInput struct {
	MatchOfferID string
	RequesterAccountID  string
}

type CreateMatchRequestUC struct {
	matchRequestRepository      matchrequest.Repository
	matchOfferRepository matchoffer.Repository
}

func NewCreateMatchRequestUC(
	matchRequestRepository matchrequest.Repository,
	matchOfferRepository matchoffer.Repository,
) *CreateMatchRequestUC {
	return &CreateMatchRequestUC{
		matchRequestRepository:      matchRequestRepository,
		matchOfferRepository: matchOfferRepository,
	}
}

func (uc *CreateMatchRequestUC) Invoke(ctx context.Context, input CreateMatchRequestInput) (*matchrequest.Entity, error) {
	// Fetch the match offer to get the owner account ID
	page, err := uc.matchOfferRepository.Find(ctx, matchoffer.DomainQuery{
		IDs: []string{input.MatchOfferID},
	})
	if err != nil {
		return nil, fmt.Errorf("error while finding match offer: %w", err)
	}
	if len(page.Entities) == 0 {
		return nil, fmt.Errorf("match offer '%s' not found", input.MatchOfferID)
	}
	offer := &page.Entities[0]

	// Prevent the owner from requesting their own offer
	if offer.OwnerAccountID == input.RequesterAccountID {
		return nil, fmt.Errorf("cannot send a match request to your own offer")
	}

	entity := matchrequest.NewMatchRequest(
		input.MatchOfferID,
		offer.OwnerAccountID,
		input.RequesterAccountID,
	)

	if err := uc.matchRequestRepository.Save(ctx, entity); err != nil {
		return nil, fmt.Errorf("error while saving match request: %w", err)
	}

	return &entity, nil
}
