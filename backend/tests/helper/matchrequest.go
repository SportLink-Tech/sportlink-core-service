package helper

import (
	"context"
	"testing"

	matchrequestuc "sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/matchrequest"
)

// MatchRequestBuilder builds and persists a matchrequest.Entity for e2e tests.
// MatchOfferID and RequesterAccountID are required; set them via WithMatchOfferID
// and WithRequesterAccountID before calling Build.
type MatchRequestBuilder struct {
	t                  *testing.T
	mrRepo             matchrequest.Repository
	moRepo             matchoffer.Repository
	matchOfferID       string
	requesterAccountID string
}

// NewMatchRequestBuilder returns a builder with no defaults for IDs,
// since both are references to entities that must exist beforehand.
func NewMatchRequestBuilder(
	t *testing.T,
	mrRepo matchrequest.Repository,
	moRepo matchoffer.Repository,
) *MatchRequestBuilder {
	t.Helper()
	return &MatchRequestBuilder{
		t:      t,
		mrRepo: mrRepo,
		moRepo: moRepo,
	}
}

func (b *MatchRequestBuilder) WithMatchOfferID(id string) *MatchRequestBuilder {
	b.matchOfferID = id
	return b
}

func (b *MatchRequestBuilder) WithRequesterAccountID(id string) *MatchRequestBuilder {
	b.requesterAccountID = id
	return b
}

// Build creates the entity, saves it via the repository and returns it.
// It calls t.Fatal on any error.
func (b *MatchRequestBuilder) Build(ctx context.Context) *matchrequest.Entity {
	b.t.Helper()

	if b.matchOfferID == "" {
		b.t.Fatal("MatchRequestBuilder: MatchOfferID is required")
	}
	if b.requesterAccountID == "" {
		b.t.Fatal("MatchRequestBuilder: RequesterAccountID is required")
	}

	uc := matchrequestuc.NewCreateMatchRequestUC(b.mrRepo, b.moRepo)
	result, err := uc.Invoke(ctx, matchrequestuc.CreateMatchRequestInput{
		MatchOfferID:       b.matchOfferID,
		RequesterAccountID: b.requesterAccountID,
	})
	if err != nil {
		b.t.Fatalf("MatchRequestBuilder: failed to save match request: %v", err)
	}

	return result
}
