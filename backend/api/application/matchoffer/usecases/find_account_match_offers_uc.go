package usecases

import (
	"context"
	"sportlink/api/domain/matchoffer"
)

type FindAccountMatchOffersResult struct {
	Entities []matchoffer.Entity
}

type FindAccountMatchOffersUC struct {
	matchOfferRepository matchoffer.Repository
}

func NewFindAccountMatchOffersUC(matchOfferRepository matchoffer.Repository) *FindAccountMatchOffersUC {
	return &FindAccountMatchOffersUC{
		matchOfferRepository: matchOfferRepository,
	}
}

func (uc *FindAccountMatchOffersUC) Invoke(ctx context.Context, query matchoffer.DomainQuery) (*FindAccountMatchOffersResult, error) {
	page, err := uc.matchOfferRepository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return &FindAccountMatchOffersResult{Entities: page.Entities}, nil
}
