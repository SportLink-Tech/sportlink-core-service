package matchrequest

import (
	"context"

	"sportlink/api/application/matchrequest/usecases"
	domain "sportlink/api/domain/matchrequest"
)

// CreateMatchRequestUseCase is implemented by the application use case; exposed for tests and mockery.
type CreateMatchRequestUseCase interface {
	Invoke(ctx context.Context, input usecases.CreateMatchRequestInput) (*domain.Entity, error)
}

// FindMatchRequestsUseCase is implemented by the application use case; exposed for tests and mockery.
type FindMatchRequestsUseCase interface {
	Invoke(ctx context.Context, ownerAccountID string) ([]domain.Entity, error)
}

// FindSentMatchRequestsUseCase is implemented by the application use case; exposed for tests and mockery.
type FindSentMatchRequestsUseCase interface {
	Invoke(ctx context.Context, requesterAccountID string, statuses []domain.Status) ([]domain.Entity, error)
}

// UpdateMatchRequestStatusUseCase is implemented by the application use case; exposed for tests and mockery.
type UpdateMatchRequestStatusUseCase interface {
	Invoke(ctx context.Context, input usecases.UpdateMatchRequestStatusInput) error
}
