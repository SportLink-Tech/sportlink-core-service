package matchannouncement_test

import (
	"context"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/matchannouncement"
	infra "sportlink/api/infrastructure/persistence/matchannouncement"
	"sportlink/dev/testcontainer"
	"sportlink/dev/utils/slice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Save(t *testing.T) {

	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := infra.NewRepository(dynamoDbClient, "SportLinkCore")

	location := domain.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	timeSlot, _ := domain.NewTimeSlot(startTime, endTime)

	testCases := []struct {
		name       string
		entity     domain.Entity
		assertions func(t *testing.T, err error)
	}{
		{
			name: "save an item successfully",
			entity: domain.NewMatchAnnouncement(
				"Thunder Strikers",
				common.Paddle,
				tomorrow,
				timeSlot,
				location,
				domain.NewSpecificCategories([]common.Category{5, 6, 7}),
				domain.StatusPending,
				time.Now().In(tz),
			),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save an item with GreaterThan category range successfully",
			entity: domain.NewMatchAnnouncement(
				"Elite Team",
				common.Tennis,
				tomorrow,
				timeSlot,
				location,
				domain.NewGreaterThanCategory(5),
				domain.StatusPending,
				time.Now().In(tz),
			),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save an item with LessThan category range successfully",
			entity: domain.NewMatchAnnouncement(
				"Beginner Team",
				common.Football,
				tomorrow,
				timeSlot,
				location,
				domain.NewLessThanCategory(3),
				domain.StatusPending,
				time.Now().In(tz),
			),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := repository.Save(ctx, testCase.entity)

			// then
			testCase.assertions(t, err)
		})
	}
}

func Test_Find(t *testing.T) {

	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := infra.NewRepository(dynamoDbClient, "SportLinkCore")

	location := domain.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	nextWeek := time.Now().In(tz).AddDate(0, 0, 7)

	testCases := []struct {
		name       string
		query      domain.DomainQuery
		on         func(t *testing.T, repository domain.Repository)
		assertions func(t *testing.T, entities []domain.Entity, err error)
	}{
		{
			name: "find announcements by sport successfully",
			query: domain.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			on: func(t *testing.T, repository domain.Repository) {
				startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
				endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
				timeSlot, _ := domain.NewTimeSlot(startTime, endTime)

				err := repository.Save(ctx, domain.NewMatchAnnouncement(
					"Paddle Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{5}),
					domain.StatusPending,
					time.Now().In(tz),
				))
				if err != nil {
					t.Fatal(err)
				}

				err = repository.Save(ctx, domain.NewMatchAnnouncement(
					"Paddle Team B",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{6}),
					domain.StatusPending,
					time.Now().In(tz),
				))
				if err != nil {
					t.Fatal(err)
				}

				// Save a Tennis announcement that should not appear
				err = repository.Save(ctx, domain.NewMatchAnnouncement(
					"Tennis Team",
					common.Tennis,
					tomorrow,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{5}),
					domain.StatusPending,
					time.Now().In(tz),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []domain.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[domain.Entity](
					entities,
					domain.Entity{TeamName: "Paddle Team A"},
					func(a, b domain.Entity) bool {
						return a.TeamName == b.TeamName
					}))
				assert.True(t, slice.Contains[domain.Entity](
					entities,
					domain.Entity{TeamName: "Paddle Team B"},
					func(a, b domain.Entity) bool {
						return a.TeamName == b.TeamName
					}))
			},
		},
		{
			name: "find announcements by multiple statuses",
			query: domain.DomainQuery{
				Statuses: []domain.Status{domain.StatusPending, domain.StatusConfirmed},
			},
			on: func(t *testing.T, repository domain.Repository) {
				startTime := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 10, 0, 0, 0, tz)
				endTime := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 12, 0, 0, 0, tz)
				timeSlot, _ := domain.NewTimeSlot(startTime, endTime)

				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Multi Status Team A",
					common.Paddle,
					nextWeek,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{2}),
					domain.StatusPending,
					time.Now().In(tz),
				))

				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Multi Status Team B",
					common.Paddle,
					nextWeek,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{2}),
					domain.StatusConfirmed,
					time.Now().In(tz),
				))

				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Multi Status Team C",
					common.Paddle,
					nextWeek,
					timeSlot,
					location,
					domain.NewSpecificCategories([]common.Category{2}),
					domain.StatusCancelled,
					time.Now().In(tz),
				))
			},
			assertions: func(t *testing.T, entities []domain.Entity, err error) {
				assert.Nil(t, err)
				assert.GreaterOrEqual(t, len(entities), 2)
				for _, entity := range entities {
					assert.Contains(t, []domain.Status{domain.StatusPending, domain.StatusConfirmed}, entity.Status)
				}
			},
		},
		{
			name: "find announcements by sport and status",
			query: domain.DomainQuery{
				Sports:   []common.Sport{common.Paddle},
				Statuses: []domain.Status{domain.StatusPending},
			},
			on: func(t *testing.T, repository domain.Repository) {
				dayAfterTomorrow := time.Now().In(tz).AddDate(0, 0, 2)
				startTime := time.Date(dayAfterTomorrow.Year(), dayAfterTomorrow.Month(), dayAfterTomorrow.Day(), 18, 0, 0, 0, tz)
				endTime := time.Date(dayAfterTomorrow.Year(), dayAfterTomorrow.Month(), dayAfterTomorrow.Day(), 20, 0, 0, 0, tz)
				timeSlot, _ := domain.NewTimeSlot(startTime, endTime)

				// Matches both filters
				betweenRange, _ := domain.NewBetweenCategories(3, 6)
				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Paddle Pending Team",
					common.Paddle,
					dayAfterTomorrow,
					timeSlot,
					location,
					betweenRange,
					domain.StatusPending,
					time.Now().In(tz),
				))

				// Wrong sport
				betweenRange2, _ := domain.NewBetweenCategories(3, 6)
				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Tennis Pending Team",
					common.Tennis,
					dayAfterTomorrow,
					timeSlot,
					location,
					betweenRange2,
					domain.StatusPending,
					time.Now().In(tz),
				))

				// Wrong status
				betweenRange3, _ := domain.NewBetweenCategories(3, 6)
				repository.Save(ctx, domain.NewMatchAnnouncement(
					"Paddle Confirmed Team",
					common.Paddle,
					dayAfterTomorrow,
					timeSlot,
					location,
					betweenRange3,
					domain.StatusConfirmed,
					time.Now().In(tz),
				))
			},
			assertions: func(t *testing.T, entities []domain.Entity, err error) {
				assert.Nil(t, err)
				assert.GreaterOrEqual(t, len(entities), 1)
				for _, entity := range entities {
					assert.Equal(t, common.Paddle, entity.Sport)
					assert.Equal(t, domain.StatusPending, entity.Status)
				}
			},
		},
		{
			name: "find announcements that do not exist",
			query: domain.DomainQuery{
				Sports: []common.Sport{common.Sport("NonExistentSport")},
			},
			on: func(t *testing.T, repository domain.Repository) {
			},
			assertions: func(t *testing.T, entities []domain.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			testCase.on(t, repository)

			// when
			page, err := repository.Find(ctx, testCase.query)

			// then
			testCase.assertions(t, page.Entities, err)
		})
	}
}
