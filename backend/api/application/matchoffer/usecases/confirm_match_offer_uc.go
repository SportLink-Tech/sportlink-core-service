package usecases

import (
	"context"
	"fmt"
	"sportlink/api/application/errors"
	"sportlink/api/domain/match"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
	"sportlink/pkg/log"
)

type ConfirmMatchOfferInput struct {
	MatchOfferID   string
	OwnerAccountID string
}

type ConfirmMatchOfferUC struct {
	matchRepository        match.Repository
	matchOfferRepository   matchoffer.Repository
	matchRequestRepository matchrequest.Repository
}

func NewConfirmMatchOfferUC(
	matchRepository match.Repository,
	matchOfferRepository matchoffer.Repository,
	matchRequestRepository matchrequest.Repository,
) *ConfirmMatchOfferUC {
	return &ConfirmMatchOfferUC{
		matchRepository:        matchRepository,
		matchOfferRepository:   matchOfferRepository,
		matchRequestRepository: matchRequestRepository,
	}
}

func (uc *ConfirmMatchOfferUC) Invoke(ctx context.Context, input ConfirmMatchOfferInput) (*match.Entity, error) {
	offer, err := uc.getMatchOffer(ctx, input.MatchOfferID)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get match offer %s", input.MatchOfferID), err)
		return nil, err
	}

	if offer.OwnerAccountID != input.OwnerAccountID {
		return nil, errors.Unauthorized("owner account ID does not match")
	}

	if !offer.IsPending() {
		err = errors.UseCaseExecutionFailed("match offer is not pending")
		log.GetLogger(ctx).Error(fmt.Sprintf("match offer %s is not pending, status: %s", input.MatchOfferID, offer.Status), err)
		return nil, err
	}

	acceptedRequests, err := uc.getAcceptedRequests(ctx, input.MatchOfferID)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get accepted requests for offer %s", input.MatchOfferID), err)
		return nil, err
	}

	newMatch := match.NewMatch(buildParticipants(input.OwnerAccountID, acceptedRequests), offer.Sport, offer.Day)

	if err = uc.matchRepository.Save(ctx, newMatch); err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to save match for offer %s", input.MatchOfferID), err)
		return nil, err
	}

	if err = uc.matchOfferRepository.Save(ctx, offer.Confirm()); err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to confirm match offer %s", input.MatchOfferID), err)
		return nil, err
	}

	if err = uc.rejectPendingRequests(ctx, input.MatchOfferID); err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to reject pending requests for offer %s", input.MatchOfferID), err)
	}

	return &newMatch, nil
}

func (uc *ConfirmMatchOfferUC) rejectPendingRequests(ctx context.Context, matchOfferID string) error {
	pending, err := uc.matchRequestRepository.Find(ctx, matchrequest.DomainQuery{
		MatchOfferIDs: []string{matchOfferID},
		Statuses:      []matchrequest.Status{matchrequest.StatusPending},
	})
	if err != nil {
		return fmt.Errorf("failed to find pending requests: %w", err)
	}
	if len(pending) == 0 {
		return nil
	}

	rejected := make([]matchrequest.Entity, len(pending))
	for i, r := range pending {
		rejected[i] = r.Reject()
	}

	return uc.matchRequestRepository.SaveAll(ctx, rejected)
}

func (uc *ConfirmMatchOfferUC) getMatchOffer(ctx context.Context, matchOfferID string) (*matchoffer.Entity, error) {
	page, err := uc.matchOfferRepository.Find(ctx, matchoffer.DomainQuery{IDs: []string{matchOfferID}})
	if err != nil {
		return nil, err
	}
	if len(page.Entities) == 0 || len(page.Entities) > 1 {
		return nil, errors.UseCaseExecutionFailed("match offer not found or multiple match offers found for the given ID")
	}
	return &page.Entities[0], nil
}

func (uc *ConfirmMatchOfferUC) getAcceptedRequests(ctx context.Context, matchOfferID string) ([]matchrequest.Entity, error) {
	return uc.matchRequestRepository.Find(ctx, matchrequest.DomainQuery{
		MatchOfferIDs: []string{matchOfferID},
		Statuses:      []matchrequest.Status{matchrequest.StatusAccepted},
	})
}

func buildParticipants(ownerAccountID string, requests []matchrequest.Entity) []string {
	participants := make([]string, 0, 1+len(requests))
	participants = append(participants, ownerAccountID)
	for _, r := range requests {
		participants = append(participants, r.RequesterAccountID)
	}
	return participants
}
