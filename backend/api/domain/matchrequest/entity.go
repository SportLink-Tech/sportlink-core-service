package matchrequest

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Entity represents a match request in the domain.
// A requester sends a request to join or challenge the owner of a match announcement.
type Entity struct {
	ID                  string
	MatchAnnouncementID string
	OwnerAccountID      string // account ID of the match announcement owner (receives the request)
	RequesterAccountID  string // account ID of the user sending the request
	Status              Status
	CreatedAt           time.Time
}

func NewMatchRequest(
	matchAnnouncementID string,
	ownerAccountID string,
	requesterAccountID string,
) Entity {
	return Entity{
		ID:                  generateMatchRequestID(),
		MatchAnnouncementID: matchAnnouncementID,
		OwnerAccountID:      ownerAccountID,
		RequesterAccountID:  requesterAccountID,
		Status:              StatusPending,
		CreatedAt:           time.Now(),
	}
}

func generateMatchRequestID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
