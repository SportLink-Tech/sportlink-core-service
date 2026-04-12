package helper

import (
	"context"
	"testing"
	"time"

	matchofferuc "sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchoffer"
)

// MatchOfferBuilder builds and persists a matchoffer.Entity for e2e tests.
// All fields have sensible defaults; override only what matters for the test.
type MatchOfferBuilder struct {
	t              *testing.T
	repo           matchoffer.Repository
	teamName       string
	sport          common.Sport
	day            time.Time
	startTime      time.Time
	endTime        time.Time
	country        string
	province       string
	locality       string
	latitude       *float64
	longitude      *float64
	categoryRange  matchoffer.CategoryRange
	status         matchoffer.Status
	ownerAccountID string
	capacity       int
}

// NewMatchOfferBuilder returns a builder with sensible defaults.
// The offer is set for tomorrow 18:00–20:00 in Buenos Aires.
func NewMatchOfferBuilder(t *testing.T, repo matchoffer.Repository) *MatchOfferBuilder {
	t.Helper()
	loc := matchoffer.NewLocation("Argentina", "Buenos Aires", "Palermo")
	tz := loc.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 18, 0, 0, 0, tz)
	end := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 20, 0, 0, 0, tz)

	return &MatchOfferBuilder{
		t:             t,
		repo:          repo,
		sport:         common.Paddle,
		day:           tomorrow,
		startTime:     start,
		endTime:       end,
		country:       "Argentina",
		province:      "Buenos Aires",
		locality:      "Palermo",
		categoryRange: matchoffer.NewGreaterThanCategory(common.Category(1)),
		status:        matchoffer.StatusPending,
	}
}

func (b *MatchOfferBuilder) WithTeamName(name string) *MatchOfferBuilder {
	b.teamName = name
	return b
}

func (b *MatchOfferBuilder) WithSport(sport common.Sport) *MatchOfferBuilder {
	b.sport = sport
	return b
}

func (b *MatchOfferBuilder) WithDay(day time.Time) *MatchOfferBuilder {
	b.day = day
	return b
}

func (b *MatchOfferBuilder) WithTimeSlot(start, end time.Time) *MatchOfferBuilder {
	b.startTime = start
	b.endTime = end
	return b
}

func (b *MatchOfferBuilder) WithLocation(country, province, locality string) *MatchOfferBuilder {
	b.country = country
	b.province = province
	b.locality = locality
	return b
}

func (b *MatchOfferBuilder) WithCoords(lat, lng float64) *MatchOfferBuilder {
	b.latitude = &lat
	b.longitude = &lng
	return b
}

func (b *MatchOfferBuilder) WithCategoryRange(cr matchoffer.CategoryRange) *MatchOfferBuilder {
	b.categoryRange = cr
	return b
}

func (b *MatchOfferBuilder) WithStatus(status matchoffer.Status) *MatchOfferBuilder {
	b.status = status
	return b
}

func (b *MatchOfferBuilder) WithOwnerAccountID(id string) *MatchOfferBuilder {
	b.ownerAccountID = id
	return b
}

func (b *MatchOfferBuilder) WithCapacity(capacity int) *MatchOfferBuilder {
	b.capacity = capacity
	return b
}

// Build creates the entity, saves it via the repository and returns it.
// It calls t.Fatal on any error.
func (b *MatchOfferBuilder) Build(ctx context.Context) *matchoffer.Entity {
	b.t.Helper()

	var location matchoffer.Location
	if b.latitude != nil && b.longitude != nil {
		location = matchoffer.NewLocationWithCoords(b.country, b.province, b.locality, *b.latitude, *b.longitude)
	} else {
		location = matchoffer.NewLocation(b.country, b.province, b.locality)
	}

	tz := location.GetTimezone()
	timeSlot, err := matchoffer.NewTimeSlot(b.startTime, b.endTime)
	if err != nil {
		b.t.Fatalf("MatchOfferBuilder: invalid time slot: %v", err)
	}

	entity := matchoffer.NewMatchOffer(
		b.teamName,
		b.sport,
		b.day.In(tz),
		timeSlot,
		location,
		b.categoryRange,
		b.status,
		time.Now().In(tz),
		b.ownerAccountID,
		b.capacity,
	)

	uc := matchofferuc.NewCreateMatchOfferUC(b.repo)
	result, err := uc.Invoke(ctx, entity)
	if err != nil {
		b.t.Fatalf("MatchOfferBuilder: failed to save match offer: %v", err)
	}

	return result
}
