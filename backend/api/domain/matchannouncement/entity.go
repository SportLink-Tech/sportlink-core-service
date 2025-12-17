package matchannouncement

import (
	"sportlink/api/domain/common"
	"time"

	"github.com/oklog/ulid/v2"
)

// Entity represents a match announcement in the domain
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
}

func NewMatchAnnouncement(
	teamName string,
	sport common.Sport,
	day time.Time,
	timeSlot TimeSlot,
	location Location,
	admittedCategories CategoryRange,
	status Status,
	createdAt time.Time,
) Entity {
	return Entity{
		ID:                 generateMatchAnnouncementID(),
		TeamName:           teamName,
		Sport:              sport,
		Day:                day,
		TimeSlot:           timeSlot,
		Location:           location,
		AdmittedCategories: admittedCategories,
		Status:             status,
		CreatedAt:          createdAt,
	}
}

// generateMatchAnnouncementID generates a ULID for the match announcement
func generateMatchAnnouncementID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
