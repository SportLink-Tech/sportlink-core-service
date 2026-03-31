package user

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Entity struct {
	ID        string
	FirstName string
	LastName  string
	PlayerIDs []string
}

// NewUser creates a new user entity with a ULID
func NewUser(firstName, lastName string, playerIDs []string) Entity {
	return Entity{
		ID:        generateUserID(),
		FirstName: firstName,
		LastName:  lastName,
		PlayerIDs: playerIDs,
	}
}

// generateUserID generates a ULID for the user
func generateUserID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
