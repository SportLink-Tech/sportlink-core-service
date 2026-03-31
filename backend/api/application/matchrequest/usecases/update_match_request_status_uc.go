package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchrequest"
)

type UpdateMatchRequestStatusInput struct {
	ID             string
	OwnerAccountID string
	NewStatus      matchrequest.Status
}

type UpdateMatchRequestStatusUC struct {
	matchRequestRepository matchrequest.Repository
}

func NewUpdateMatchRequestStatusUC(matchRequestRepository matchrequest.Repository) *UpdateMatchRequestStatusUC {
	return &UpdateMatchRequestStatusUC{
		matchRequestRepository: matchRequestRepository,
	}
}

func (uc *UpdateMatchRequestStatusUC) Invoke(ctx context.Context, input UpdateMatchRequestStatusInput) error {
	err := uc.matchRequestRepository.UpdateStatus(ctx, input.ID, input.OwnerAccountID, input.NewStatus)
	if err != nil {
		return fmt.Errorf("error while updating match request status: %w", err)
	}
	return nil
}
