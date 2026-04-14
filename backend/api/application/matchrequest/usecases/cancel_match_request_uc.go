package usecases

import (
	"context"
	"fmt"
	"sportlink/api/application/errors"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
	"sportlink/pkg/log"
)

type CancelMatchRequestInput struct {
	MatchRequestId     string
	RequesterAccountID string
}

type CancelMatchRequestUC struct {
	matchRequestRepository matchrequest.Repository
	matchOfferRepository   matchoffer.Repository
}

func NewCancelMatchRequestUC(
	matchRequestRepository matchrequest.Repository,
	matchOfferRepository matchoffer.Repository,
) *CancelMatchRequestUC {
	return &CancelMatchRequestUC{
		matchRequestRepository: matchRequestRepository,
		matchOfferRepository:   matchOfferRepository,
	}
}

func (uc *CancelMatchRequestUC) Invoke(ctx context.Context, input CancelMatchRequestInput) (*matchrequest.Entity, error) {
	matchReq, err := getMatchRequest(ctx, uc.matchRequestRepository, input.MatchRequestId)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get match request %s", input.MatchRequestId), err)
		return nil, err
	}

	if matchReq.IsRejected() {
		log.GetLogger(ctx).Error(fmt.Sprintf("match request %s is already rejected, status: %s", input.MatchRequestId, matchReq.Status), err)
		return nil, errors.UseCaseExecutionFailed("match request is already rejected")
	}

	matchOffer, err := getMatchOffer(ctx, uc.matchOfferRepository, matchReq.MatchOfferID)
	if err != nil {
		log.GetLogger(ctx).Error(fmt.Sprintf("failed to get match offer %s", matchReq.MatchOfferID), err)
		return nil, err
	}

	if matchOffer.IsConfirm() {
		log.GetLogger(ctx).Error(fmt.Sprintf("match offer %s is already confirmed, status: %s", matchReq.MatchOfferID, matchOffer.Status), err)
		return nil, errors.UseCaseExecutionFailed("match offer is already confirmed")
	}

	canceled := matchReq.Cancel()
	if err = uc.matchRequestRepository.Save(ctx, canceled); err != nil {
		return nil, fmt.Errorf("error while cancelling match request: %w", err)
	}
	return &canceled, nil
}
