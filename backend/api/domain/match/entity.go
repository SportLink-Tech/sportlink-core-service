package match

import (
	"sportlink/api/domain/common"
	"time"

	"github.com/oklog/ulid/v2"
)

type Entity struct {
	ID             string
	MatchOfferID   string
	Participants   []string
	Sport          common.Sport
	Day            time.Time
	Status         Status
	Result         *Result
	WinnerAccountID string
	CreatedAt      time.Time
}

func NewMatch(
	matchOfferID string,
	participants []string,
	sport common.Sport,
	day time.Time,
) Entity {
	return Entity{
		ID:           generateMatchID(),
		MatchOfferID: matchOfferID,
		Participants: participants,
		Sport:        sport,
		Day:          day,
		Status:       StatusAccepted,
		CreatedAt:    time.Now().UTC(),
	}
}

func generateMatchID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
