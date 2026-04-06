package usecases

import (
	"context"
	"sportlink/api/domain/matchoffer"
)

type DeleteMatchOfferUC struct {
	matchOfferRepository matchoffer.Repository
}

func NewDeleteMatchOfferUC(repo matchoffer.Repository) *DeleteMatchOfferUC {
	return &DeleteMatchOfferUC{matchOfferRepository: repo}
}

func (uc *DeleteMatchOfferUC) Invoke(ctx context.Context, offerID string) error {
	return uc.matchOfferRepository.Delete(ctx, offerID)
}
