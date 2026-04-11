package match

import (
	"sportlink/api/domain/common"
	"time"

	"github.com/oklog/ulid/v2"
)

type Entity struct {
	ID               string
	LocalAccountID   string
	VisitorAccountID string
	Sport            common.Sport
	Day              time.Time
	Status           Status
	Result           *Result
	WinnerAccountID  string
	CreatedAt        time.Time
}

func NewMatch(
	localAccountID string,
	visitorAccountID string,
	sport common.Sport,
	day time.Time,
) Entity {
	return Entity{
		ID:               generateMatchID(),
		LocalAccountID:   localAccountID,
		VisitorAccountID: visitorAccountID,
		Sport:            sport,
		Day:              day,
		Status:           StatusAccepted,
		CreatedAt:        time.Now().UTC(),
	}
}

func generateMatchID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
