package matchoffer

import (
	"sportlink/api/domain/common"
	"time"

	"github.com/oklog/ulid/v2"
)

// Entity represents a match offer in the domain
// ID is generated automatically when the entity is created
// ExpiresAt is the responsibility of the persistence layer
type Entity struct {
	ID                 string
	TeamName           string
	Sport              common.Sport
	Day                time.Time
	TimeSlot           TimeSlot
	Location           Location
	AdmittedCategories CategoryRange
	Status             Status
	CreatedAt          time.Time
	OwnerAccountID     string
	Capacity           int // 0 = no auto-confirm; >0 = total spots (owner + accepted requesters)
}

func NewMatchOffer(
	teamName string,
	sport common.Sport,
	day time.Time,
	timeSlot TimeSlot,
	location Location,
	admittedCategories CategoryRange,
	status Status,
	createdAt time.Time,
	ownerAccountID string,
	capacity int,
) Entity {
	return Entity{
		ID:                 generateMatchOfferID(),
		TeamName:           teamName,
		Sport:              sport,
		Day:                day,
		TimeSlot:           timeSlot,
		Location:           location,
		AdmittedCategories: admittedCategories,
		Status:             status,
		CreatedAt:          createdAt,
		OwnerAccountID:     ownerAccountID,
		Capacity:           capacity,
	}
}

func (s Entity) Confirm() Entity {
	s.Status = StatusConfirmed
	return s
}

func (s Entity) IsPending() bool {
	loc, _ := time.LoadLocation("America/New_York")
	if s.IsExpire(loc) {
		return false
	}
	return s.Status == StatusPending

}

// IsExpire returns true if the match has already ended based on the current time.
// It uses TimeSlot.EndTime as the expiry boundary — once the match end time has
// passed, the offer is considered expired regardless of its Status.
// loc defaults to GMT-3 (Argentina) when nil.
func (s Entity) IsExpire(loc *time.Location) bool {
	if loc == nil {
		loc = time.FixedZone("GMT-3", -3*60*60)
	}
	return time.Now().In(loc).After(s.TimeSlot.EndTime.In(loc))
}

// generateMatchOfferID generates a ULID for the match offer
func generateMatchOfferID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
