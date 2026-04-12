package match

import "context"

type DomainQuery struct {
	AccountID string
	Statuses  []Status
}

type Repository interface {
	// Save persists a match. Writes two records (one per account) so both
	// All participants can efficiently list their matches.
	Save(ctx context.Context, entity Entity) error

	// Find returns all matches for the given account (as local or visitor).
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)

	// FindByID returns a single match by ID, scoped to one of its participant accounts.
	FindByID(ctx context.Context, accountID, matchID string) (*Entity, error)
}
