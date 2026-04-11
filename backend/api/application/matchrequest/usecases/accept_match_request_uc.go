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

type AcceptMatchRequestInput struct {
	MatchRequestId string
	OwnerAccountID string
}

type AcceptMatchRequestUC struct {
	matchRepository        match.Repository
	matchRequestRepository matchrequest.Repository
	matchOfferRepository   matchoffer.Repository
}

func NewAcceptMatchRequestUC(
	matchRepository match.Repository,
	matchRequestRepository matchrequest.Repository,
	matchOfferRepository matchoffer.Repository,
) *AcceptMatchRequestUC {
	return &AcceptMatchRequestUC{
		matchRepository:        matchRepository,
		matchRequestRepository: matchRequestRepository,
		matchOfferRepository:   matchOfferRepository,
	}
}

func (uc *AcceptMatchRequestUC) Invoke(ctx context.Context, input AcceptMatchRequestInput) (*match.Entity, error) {
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

	newMatch := match.NewMatch(matchReq.OwnerAccountID, matchReq.RequesterAccountID, matchOffer.Sport, matchOffer.Day)
	err = uc.matchRepository.Save(ctx, newMatch)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to save match %s", newMatch.ID), err)
		return nil, err
	}

	err = uc.matchRequestRepository.Save(ctx, matchReq.Accept())
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to save accepted match request %s", input.MatchRequestId), err)
		return nil, err
	}

	err = uc.matchOfferRepository.Save(ctx, matchOffer.Confirm())
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to save confirmed match offer %s", matchReq.MatchOfferID), err)
		return nil, err
	}

	return &newMatch, nil
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
