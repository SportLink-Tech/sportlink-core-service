package usecases

import (
	"context"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
	"sportlink/pkg/slices"
)

// SearchMatchOffersInput contains the search criteria and the account performing the search.
type SearchMatchOffersInput struct {
	ViewerAccountID string
	Query           matchoffer.DomainQuery
}

// SearchMatchOffersUC returns match offers available for a given account, automatically
// excluding offers owned by that account and offers where the account already has a
// PENDING or ACCEPTED match request.
type SearchMatchOffersUC struct {
	matchOfferRepo   matchoffer.Repository
	matchRequestRepo matchrequest.Repository
}

func NewSearchMatchOffersUC(
	matchOfferRepo matchoffer.Repository,
	matchRequestRepo matchrequest.Repository,
) *SearchMatchOffersUC {
	return &SearchMatchOffersUC{
		matchOfferRepo:   matchOfferRepo,
		matchRequestRepo: matchRequestRepo,
	}
}

func (uc *SearchMatchOffersUC) Invoke(ctx context.Context, input SearchMatchOffersInput) (*FindMatchOfferResult, error) {
	requestedOfferIDs, err := uc.fetchOfferIDsWithActiveRequestByViewer(ctx, input.ViewerAccountID)
	if err != nil {
		return nil, err
	}

	page, err := uc.matchOfferRepo.Find(ctx, input.Query)
	if err != nil {
		return nil, err
	}

	available := uc.excludeOffersUnavailableForViewer(page.Entities, input.ViewerAccountID, requestedOfferIDs)

	pageInfo := CalculatePageInfo(
		input.Query.Limit,
		input.Query.Offset,
		adjustTotal(page.Total, page.Entities, available),
	)

	return &FindMatchOfferResult{
		Entities: available,
		Page:     pageInfo,
	}, nil
}

// fetchOfferIDsWithActiveRequestByViewer returns the IDs of match offers for which the viewer
// already has a PENDING or ACCEPTED match request.
func (uc *SearchMatchOffersUC) fetchOfferIDsWithActiveRequestByViewer(
	ctx context.Context,
	viewerAccountID string,
) ([]string, error) {
	requests, err := uc.matchRequestRepo.Find(ctx, matchrequest.DomainQuery{
		RequesterAccountIDs: []string{viewerAccountID},
		Statuses:            []matchrequest.Status{matchrequest.StatusPending, matchrequest.StatusAccepted},
	})
	if err != nil {
		return nil, err
	}
	return slices.Map(requests, func(r matchrequest.Entity) string {
		return r.MatchOfferID
	}), nil
}

// excludeOffersUnavailableForViewer filters out offers owned by the viewer and offers
// for which the viewer already has an active match request.
func (uc *SearchMatchOffersUC) excludeOffersUnavailableForViewer(
	offers []matchoffer.Entity,
	viewerAccountID string,
	requestedOfferIDs []string,
) []matchoffer.Entity {
	return slices.Filter(offers, func(e matchoffer.Entity) bool {
		return e.OwnerAccountID != viewerAccountID && !slices.Contains(requestedOfferIDs, e.ID)
	})
}

// adjustTotal corrects the total count by subtracting the number of offers excluded on the current page.
// This is a best-effort approximation since DynamoDB's total includes offers that were filtered out.
func adjustTotal(total int, original, filtered []matchoffer.Entity) int {
	adjusted := total - (len(original) - len(filtered))
	if adjusted < 0 {
		return 0
	}
	return adjusted
}
