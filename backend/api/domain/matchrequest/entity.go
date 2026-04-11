package matchrequest

import (
	"fmt"
	"time"
)

// Entity represents a match request in the domain.
// A requester sends a request to join or challenge the owner of a match announcement.
type Entity struct {
	ID                 string
	MatchOfferID       string
	OwnerAccountID     string // account ID of the match announcement owner (receives the request)
	RequesterAccountID string // account ID of the user sending the request
	Status             Status
	CreatedAt          time.Time
}

func NewMatchRequest(
	matchOfferID string,
	ownerAccountID string,
	requesterAccountID string,
) Entity {
	return Entity{
		ID:                 GenerateMatchRequestID(requesterAccountID, matchOfferID),
		MatchOfferID:       matchOfferID,
		OwnerAccountID:     ownerAccountID,
		RequesterAccountID: requesterAccountID,
		Status:             StatusPending,
		CreatedAt:          time.Now(),
	}
}
func (s Entity) Accept() Entity {
	s.Status = StatusAccepted
	return s
}

func (s Entity) IsPending() bool {
	return s.Status == StatusPending
}

func (s Entity) IsRejected() bool {
	return s.Status == StatusRejected
}

func (s Entity) IsAccepted() bool {
	return s.Status == StatusAccepted
}

// GenerateMatchRequestID returns the composite sort key for a match request.
// Format: AccountId#<requesterAccountID>#MatchOfferId#<matchOfferID>
func GenerateMatchRequestID(requesterAccountID, matchOfferID string) string {
	return fmt.Sprintf("AccountId#%s#MatchOfferId#%s", requesterAccountID, matchOfferID)
}

// GenerateMatchRequestIDPrefix returns the prefix used to query all requests from a requester.
func GenerateMatchRequestIDPrefix(requesterAccountID string) string {
	return fmt.Sprintf("AccountId#%s#MatchOfferId#", requesterAccountID)
}
