package player

import (
	"sportlink/api/domain/common"
	"time"

	"github.com/oklog/ulid/v2"
)

type Entity struct {
	ID       string
	Category common.Category
	Sport    common.Sport
}

// NewPlayer creates a new player entity with a ULID
func NewPlayer(category common.Category, sport common.Sport) Entity {
	return Entity{
		ID:       generatePlayerID(),
		Category: category,
		Sport:    sport,
	}
}

// generatePlayerID generates a ULID for the player
func generatePlayerID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
