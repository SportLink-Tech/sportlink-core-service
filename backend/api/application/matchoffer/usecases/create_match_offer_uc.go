package usecases

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchoffer"
	"time"
)

type CreateMatchOfferUC struct {
	matchOfferRepository matchoffer.Repository
}

func NewCreateMatchOfferUC(
	matchOfferRepository matchoffer.Repository,
) *CreateMatchOfferUC {
	return &CreateMatchOfferUC{
		matchOfferRepository: matchOfferRepository,
	}
}

func (uc *CreateMatchOfferUC) Invoke(ctx context.Context, input matchoffer.Entity) (*matchoffer.Entity, error) {
	if err := uc.validateOffer(input); err != nil {
		return nil, err
	}

	if err := uc.matchOfferRepository.Save(ctx, input); err != nil {
		return nil, fmt.Errorf("error while inserting match offer in database: %w", err)
	}

	return &input, nil
}

func (uc *CreateMatchOfferUC) validateOffer(input matchoffer.Entity) error {
	if input.Sport == "" {
		return fmt.Errorf("sport cannot be empty")
	}

	now := time.Now().In(input.Location.GetTimezone())
	if input.Day.Before(now.Truncate(24 * time.Hour)) {
		return fmt.Errorf("day cannot be in the past")
	}

	if input.TimeSlot.StartTime.IsZero() || input.TimeSlot.EndTime.IsZero() {
		return fmt.Errorf("time slot cannot be empty")
	}

	if input.TimeSlot.EndTime.Before(input.TimeSlot.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}

	if input.Location.Country == "" || input.Location.Province == "" || input.Location.Locality == "" {
		return fmt.Errorf("location must have country, province and locality")
	}

	if !input.Status.IsValid() {
		return fmt.Errorf("invalid status")
	}

	if input.CreatedAt.IsZero() {
		return fmt.Errorf("created at cannot be empty")
	}

	return nil
}
