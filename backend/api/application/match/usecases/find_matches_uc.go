package usecases

import (
	"context"
	"sportlink/api/domain/match"
	"sportlink/api/domain/matchoffer"
)

type FindMatchesInput struct {
	AccountID string
	Statuses  []match.Status
}

type MatchWithOffer struct {
	Match match.Entity
	Offer *matchoffer.Entity
}

type FindMatchesUC struct {
	matchRepository     match.Repository
	matchOfferRepository matchoffer.Repository
}

func NewFindMatchesUC(matchRepository match.Repository, matchOfferRepository matchoffer.Repository) *FindMatchesUC {
	return &FindMatchesUC{
		matchRepository:      matchRepository,
		matchOfferRepository: matchOfferRepository,
	}
}

func (uc *FindMatchesUC) Invoke(ctx context.Context, input FindMatchesInput) (*[]MatchWithOffer, error) {
	entities, err := uc.matchRepository.Find(ctx, match.DomainQuery{
		AccountID: input.AccountID,
		Statuses:  input.Statuses,
	})
	if err != nil {
		return nil, err
	}

	offerMap, err := uc.resolveOffers(ctx, entities)
	if err != nil {
		return nil, err
	}

	results := make([]MatchWithOffer, len(entities))
	for i, e := range entities {
		offer := offerMap[e.MatchOfferID]
		results[i] = MatchWithOffer{Match: e, Offer: offer}
	}
	return &results, nil
}

func (uc *FindMatchesUC) resolveOffers(ctx context.Context, entities []match.Entity) (map[string]*matchoffer.Entity, error) {
	offerIDSet := make(map[string]struct{})
	for _, e := range entities {
		if e.MatchOfferID != "" {
			offerIDSet[e.MatchOfferID] = struct{}{}
		}
	}
	if len(offerIDSet) == 0 {
		return map[string]*matchoffer.Entity{}, nil
	}

	ids := make([]string, 0, len(offerIDSet))
	for id := range offerIDSet {
		ids = append(ids, id)
	}

	page, err := uc.matchOfferRepository.Find(ctx, matchoffer.DomainQuery{IDs: ids})
	if err != nil {
		return nil, err
	}

	offerMap := make(map[string]*matchoffer.Entity, len(page.Entities))
	for i := range page.Entities {
		o := page.Entities[i]
		offerMap[o.ID] = &o
	}
	return offerMap, nil
}
