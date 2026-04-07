package usecases

import (
	"context"
	"sportlink/api/domain/matchoffer"
)

type RetrieveMatchOfferUC struct {
	matchOfferRepository matchoffer.Repository
}

func NewRetrieveMatchOfferUC(repo matchoffer.Repository) *RetrieveMatchOfferUC {
	return &RetrieveMatchOfferUC{matchOfferRepository: repo}
}

func (rm *RetrieveMatchOfferUC) Invoke(ctx context.Context, query matchoffer.DomainQuery) (*matchoffer.Entity, error) {
	result, err := rm.matchOfferRepository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return &result.Entities[0], nil
}
