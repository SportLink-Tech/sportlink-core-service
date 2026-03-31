package matchrequest

import "context"

// DomainQuery represents the search criteria for match requests
type DomainQuery struct {
	IDs                 []string // Search by specific request IDs
	OwnerAccountIDs     []string // Search by match announcement owner account IDs
	RequesterAccountIDs []string // Search by requester account IDs
	Statuses            []Status // Search by statuses
}

// Repository defines the persistence operations for match requests
type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)
	UpdateStatus(ctx context.Context, id string, ownerAccountID string, newStatus Status) error
}
