package usecases

import (
	"context"
	"fmt"
	"sportlink/api/application/errors"
	appevents "sportlink/api/application/events"
	matchofferevent "sportlink/api/application/matchoffer/events"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
	"sportlink/pkg/log"
)

type AcceptMatchRequestInput struct {
	MatchRequestId string
	OwnerAccountID string
}

type AcceptMatchRequestUC struct {
	matchRequestRepository matchrequest.Repository
	matchOfferRepository   matchoffer.Repository
	publisher              appevents.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]
}

func NewAcceptMatchRequestUC(
	matchRequestRepository matchrequest.Repository,
	matchOfferRepository matchoffer.Repository,
	publisher appevents.Publisher[matchofferevent.MatchOfferCapacityReachedEvent],
) *AcceptMatchRequestUC {
	return &AcceptMatchRequestUC{
		matchRequestRepository: matchRequestRepository,
		matchOfferRepository:   matchOfferRepository,
		publisher:              publisher,
	}
}

func (uc *AcceptMatchRequestUC) Invoke(ctx context.Context, input AcceptMatchRequestInput) (*matchrequest.Entity, error) {
	matchReq, err := uc.getMatchRequest(ctx, input.MatchRequestId)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get match request %s", input.MatchRequestId), err)
		return nil, err
	}

	if matchReq.OwnerAccountID != input.OwnerAccountID {
		return nil, errors.Unauthorized("owner account ID does not match")
	}

	if !matchReq.IsPending() {
		err = errors.UseCaseExecutionFailed("match request is not pending")
		log.GetLogger(ctx).Error(fmt.Sprintf("match request %s is not pending, status: %s", input.MatchRequestId, matchReq.Status), err)
		return nil, err
	}

	matchOffer, err := uc.getMatchOffer(ctx, matchReq.MatchOfferID)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get match offer %s", matchReq.MatchOfferID), err)
		return nil, err
	}
	if !matchOffer.IsPending() {
		err = errors.UseCaseExecutionFailed("match offer is not pending")
		log.GetLogger(ctx).Error(fmt.Sprintf("match offer %s is not pending, status: %s", matchReq.MatchOfferID, matchOffer.Status), err)
		return nil, err
	}

	accepted := matchReq.Accept()
	if err = uc.matchRequestRepository.Save(ctx, accepted); err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to save accepted match request %s", input.MatchRequestId), err)
		return nil, err
	}

	if err = uc.tryPublishIfCapacityReached(ctx, matchOffer); err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to publish capacity reached event for offer %s", matchOffer.ID), err)
	}

	return &accepted, nil
}

func (uc *AcceptMatchRequestUC) tryPublishIfCapacityReached(ctx context.Context, offer *matchoffer.Entity) error {
	if offer.Capacity == 0 {
		return nil
	}
	count, err := uc.countAcceptedRequests(ctx, offer.ID)
	if err != nil {
		return err
	}
	if count == offer.Capacity-1 {
		return uc.publisher.Publish(ctx, matchofferevent.MatchOfferCapacityReachedEvent{
			MatchOfferID:   offer.ID,
			OwnerAccountID: offer.OwnerAccountID,
		})
	}
	return nil
}

func (uc *AcceptMatchRequestUC) countAcceptedRequests(ctx context.Context, matchOfferID string) (int, error) {
	requests, err := uc.matchRequestRepository.Find(ctx, matchrequest.DomainQuery{
		MatchOfferIDs: []string{matchOfferID},
		Statuses:      []matchrequest.Status{matchrequest.StatusAccepted},
	})
	if err != nil {
		return 0, err
	}
	return len(requests), nil
}

func (uc *AcceptMatchRequestUC) getMatchRequest(ctx context.Context, matchRequestId string) (*matchrequest.Entity, error) {
	matchReqs, err := uc.matchRequestRepository.Find(ctx, matchrequest.DomainQuery{IDs: []string{matchRequestId}})
	if err != nil {
		return nil, err
	}
	if len(matchReqs) == 0 || len(matchReqs) > 1 {
		return nil, errors.UseCaseExecutionFailed("match request not found or multiple match requests found for the given ID")
	}
	return &matchReqs[0], nil
}

func (uc *AcceptMatchRequestUC) getMatchOffer(ctx context.Context, matchOfferID string) (*matchoffer.Entity, error) {
	matchOffersPage, err := uc.matchOfferRepository.Find(ctx, matchoffer.DomainQuery{IDs: []string{matchOfferID}})
	if err != nil {
		return nil, err
	}
	if len(matchOffersPage.Entities) == 0 || len(matchOffersPage.Entities) > 1 {
		return nil, errors.UseCaseExecutionFailed("match offer not found or multiple match offers found for the given ID")
	}
	return &matchOffersPage.Entities[0], nil
}
