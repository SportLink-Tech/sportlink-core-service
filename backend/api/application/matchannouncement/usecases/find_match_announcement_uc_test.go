package usecases_test

import (
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
		then  func(t *testing.T, result *[]matchannouncement.Entity, err error)
	}{
		{
			name: "find announcements successfully - multiple results by sport",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 2)
				assert.Equal(t, "Thunder Strikers", (*result)[0].TeamName)
				assert.Equal(t, "Elite Padel Team", (*result)[1].TeamName)
				assert.Equal(t, common.Paddle, (*result)[0].Sport)
			},
		},
		{
			name: "find announcements successfully - single result by sport and status",
			query: matchannouncement.DomainQuery{
				Sports:   []common.Sport{common.Tennis},
				Statuses: []matchannouncement.Status{matchannouncement.StatusConfirmed},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Tennis &&
						len(query.Statuses) == 1 &&
						query.Statuses[0] == matchannouncement.StatusConfirmed
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "Tennis Pros", (*result)[0].TeamName)
				assert.Equal(t, common.Tennis, (*result)[0].Sport)
				assert.Equal(t, matchannouncement.StatusConfirmed, (*result)[0].Status)
			},
		},
		{
			name: "find announcements by multiple statuses successfully",
			query: matchannouncement.DomainQuery{
				Statuses: []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Statuses) == 2 &&
						reflect.DeepEqual(query.Statuses, []matchannouncement.Status{
							matchannouncement.StatusPending,
							matchannouncement.StatusConfirmed,
						})
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 2)
				assert.Contains(t, []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				}, (*result)[0].Status)
				assert.Contains(t, []matchannouncement.Status{
					matchannouncement.StatusPending,
					matchannouncement.StatusConfirmed,
				}, (*result)[1].Status)
			},
		},
		{
			name: "find announcements by date range successfully",
			query: matchannouncement.DomainQuery{
				FromDate: time.Now().In(tz),
				ToDate:   time.Now().In(tz).AddDate(0, 0, 10),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return !query.FromDate.IsZero() && !query.ToDate.IsZero()
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "Future Match Team", (*result)[0].TeamName)
			},
		},
		{
			name: "find announcements by location successfully",
			query: matchannouncement.DomainQuery{
				Location: &matchannouncement.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "CABA",
				},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return query.Location != nil &&
						query.Location.Country == "Argentina" &&
						query.Location.Locality == "CABA"
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "CABA", (*result)[0].Location.Locality)
			},
		},
		{
			name: "find announcements with multiple filters - sport and status",
			query: matchannouncement.DomainQuery{
				Sports:   []common.Sport{common.Paddle},
				Statuses: []matchannouncement.Status{matchannouncement.StatusPending},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						len(query.Statuses) == 1 &&
						query.Statuses[0] == matchannouncement.StatusPending
				})).Return([]matchannouncement.Entity{
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
				}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, common.Paddle, (*result)[0].Sport)
				assert.Equal(t, matchannouncement.StatusPending, (*result)[0].Status)
			},
		},
		{
			name: "find announcements successfully - no results",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Sport("NonExistentSport")},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1
				})).Return([]matchannouncement.Entity{}, nil)
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Empty(t, *result)
			},
		},
		{
			name: "find announcements fails - repository error",
			query: matchannouncement.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query matchannouncement.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return([]matchannouncement.Entity{}, fmt.Errorf("database connection error"))
			},
			then: func(t *testing.T, result *[]matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database connection error", err.Error())
				assert.Nil(t, result)
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
			result, err := uc.Invoke(tt.query)

			// then
			tt.then(t, result, err)
		})
	}
}
