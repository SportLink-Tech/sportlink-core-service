package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"sportlink/api/application/matchannouncement/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	mmocks "sportlink/mocks/api/domain/matchannouncement"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindMatchAnnouncementUC_Invoke(t *testing.T) {

	location := matchannouncement.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	nextWeek := time.Now().In(tz).AddDate(0, 0, 7)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	timeSlot, _ := matchannouncement.NewTimeSlot(startTime, endTime)

	tests := []struct {
		name  string
		query matchannouncement.DomainQuery
		on    func(t *testing.T, repository *mmocks.Repository)
		then  func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error)
	}{
		{
			name: "given valid query when finding announcements then returns paginated results",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Limit:  9,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "Thunder Strikers",
							Sport:              common.Paddle,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{5, 6, 7}),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
						{
							TeamName:           "Elite Padel Team",
							Sport:              common.Paddle,
							Day:                nextWeek,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewGreaterThanCategory(5),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 25,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 2)
				assert.Equal(t, "Thunder Strikers", result.Entities[0].TeamName)
				assert.Equal(t, "Elite Padel Team", result.Entities[1].TeamName)
				assert.Equal(t, common.Paddle, result.Entities[0].Sport)
				assert.Equal(t, 25, result.Page.Total)
				assert.Equal(t, 1, result.Page.Number)
				assert.Equal(t, 3, result.Page.OutOf)
			},
		},
		{
			name: "given second page when finding announcements then returns correct page number",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Tennis},
				Limit:  9,
				Offset: 9,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Tennis
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "Tennis Pros",
							Sport:              common.Tennis,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{4, 5}),
							Status:             matchannouncement.StatusConfirmed,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 20,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, "Tennis Pros", result.Entities[0].TeamName)
				assert.Equal(t, common.Tennis, result.Entities[0].Sport)
				assert.Equal(t, matchannouncement.StatusConfirmed, result.Entities[0].Status)
				assert.Equal(t, 20, result.Page.Total)
				assert.Equal(t, 2, result.Page.Number)
				assert.Equal(t, 3, result.Page.OutOf)
			},
		},
		{
			name: "given query with multiple statuses when finding announcements then returns matching results",
			query: matchannouncement.DomainQuery{
				Statuses: []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				},
				Limit:  9,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Statuses) == 2 &&
						reflect.DeepEqual(query.Statuses, []matchannouncement.Status{
							matchannouncement.StatusPending,
							matchannouncement.StatusConfirmed,
						})
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "Team Pending",
							Sport:              common.Football,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{1}),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
						{
							TeamName:           "Team Confirmed",
							Sport:              common.Football,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{1}),
							Status:             matchannouncement.StatusConfirmed,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 2,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 2)
				assert.Contains(t, []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				}, result.Entities[0].Status)
				assert.Contains(t, []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				}, result.Entities[1].Status)
				assert.Equal(t, 2, result.Page.Total)
				assert.Equal(t, 1, result.Page.Number)
				assert.Equal(t, 1, result.Page.OutOf)
			},
		},
		{
			name: "given query with date range when finding announcements then returns matching results",
			query: matchannouncement.DomainQuery{
				FromDate: time.Now().In(tz),
				ToDate:   time.Now().In(tz).AddDate(0, 0, 10),
				Limit:    9,
				Offset:   0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return !query.FromDate.IsZero() && !query.ToDate.IsZero()
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "Future Match Team",
							Sport:              common.Paddle,
							Day:                nextWeek,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{3, 4}),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 1,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, "Future Match Team", result.Entities[0].TeamName)
				assert.Equal(t, 1, result.Page.Total)
			},
		},
		{
			name: "given query with location when finding announcements then returns matching results",
			query: matchannouncement.DomainQuery{
				Location: &matchannouncement.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "CABA",
				},
				Limit:  9,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return query.Location != nil &&
						query.Location.Country == "Argentina" &&
						query.Location.Locality == "CABA"
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "CABA Local Team",
							Sport:              common.Football,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewLessThanCategory(3),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 1,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, "CABA", result.Entities[0].Location.Locality)
			},
		},
		{
			name: "given query with multiple filters when finding announcements then returns matching results",
			query: matchannouncement.DomainQuery{
				Sports:   []common.Sport{common.Paddle},
				Statuses: []matchannouncement.Status{matchannouncement.StatusPending},
				Limit:    9,
				Offset:   0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						len(query.Statuses) == 1 &&
						query.Statuses[0] == matchannouncement.StatusPending
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "Paddle Seekers",
							Sport:              common.Paddle,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{5, 6}),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 1,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, common.Paddle, result.Entities[0].Sport)
				assert.Equal(t, matchannouncement.StatusPending, result.Entities[0].Status)
			},
		},
		{
			name: "given query with no results when finding announcements then returns empty result",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Sport("NonExistentSport")},
				Limit:  9,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{},
					Total:    0,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Empty(t, result.Entities)
				assert.Equal(t, 0, result.Page.Total)
				assert.Equal(t, 0, result.Page.OutOf)
			},
		},
		{
			name: "given repository error when finding announcements then returns error",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Limit:  9,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return(matchannouncement.Page{}, fmt.Errorf("database connection error"))
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database connection error", err.Error())
				assert.Nil(t, result)
			},
		},
		{
			name: "given query without limit when finding announcements then returns single page",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Limit:  0,
				Offset: 0,
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{
						{
							TeamName:           "All Results Team",
							Sport:              common.Paddle,
							Day:                tomorrow,
							TimeSlot:           timeSlot,
							Location:           location,
							AdmittedCategories: matchannouncement.NewSpecificCategories([]common.Category{5}),
							Status:             matchannouncement.StatusPending,
							CreatedAt:          time.Now().In(tz),
						},
					},
					Total: 1,
				}, nil)
			},
			then: func(t *testing.T, result *usecases.FindMatchAnnouncementResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, 1, result.Page.Number)
				assert.Equal(t, 1, result.Page.OutOf)
				assert.Equal(t, 1, result.Page.Total)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			//set up
			repo := &mmocks.Repository{}
			uc := usecases.NewFindMatchAnnouncementUC(repo)

			// given
			tt.on(t, repo)

			// when
			result, err := uc.Invoke(context.Background(), tt.query)

			// then
			tt.then(t, result, err)
		})
	}
}
