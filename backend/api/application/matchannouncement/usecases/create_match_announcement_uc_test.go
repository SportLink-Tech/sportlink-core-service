package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"sportlink/api/application/matchannouncement/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/domain/team"
	mmocks "sportlink/mocks/api/domain/matchannouncement"
	tmocks "sportlink/mocks/api/domain/team"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCreateMatchAnnouncementUC(t *testing.T) {

	location := matchannouncement.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	yesterday := time.Now().In(tz).AddDate(0, 0, -1)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 10, 0, 0, 0, tz)
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 12, 0, 0, 0, tz)

	timeSlot, _ := matchannouncement.NewTimeSlot(startTime, endTime)
	pastTimeSlot, _ := matchannouncement.NewTimeSlot(yesterdayStart, yesterdayEnd)
	categoryRange := matchannouncement.NewSpecificCategories([]common.Category{5, 6, 7})
	greaterThanRange := matchannouncement.NewGreaterThanCategory(5)

	tests := []struct {
		name  string
		input matchannouncement.Entity
		given func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository)
		then  func(t *testing.T, result *matchannouncement.Entity, err error)
	}{
		{
			name: "given valid announcement and team exists when invoke then saves and returns entity",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save announcement
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity matchannouncement.Entity) bool {
					return entity.TeamName == "Thunder Strikers" &&
						entity.Sport == common.Paddle &&
						entity.Status == matchannouncement.StatusPending
				})).Return(nil)
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "Thunder Strikers", result.TeamName)
				assert.Equal(t, common.Paddle, result.Sport)
			},
		},
		{
			name: "given greater than category range and team exists when invoke then saves and returns entity",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: greaterThanRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save announcement
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity matchannouncement.Entity) bool {
					return entity.TeamName == "Thunder Strikers" &&
						entity.AdmittedCategories.Type == matchannouncement.RangeTypeGreaterThan
				})).Return(nil)
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, matchannouncement.RangeTypeGreaterThan, result.AdmittedCategories.Type)
			},
		},
		{
			name: "given team does not exist when invoke then returns error and does not save",
			input: matchannouncement.Entity{
				TeamName:           "NonExistent Team",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team does not exist (empty slice)
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "NonExistent Team" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{}, nil)
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "team 'NonExistent Team' for sport 'Paddle' does not exist")
			},
		},
		{
			name: "given team repository fails when invoke then returns wrapped find error",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team repository error
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers"
				})).Return([]team.Entity{}, fmt.Errorf("database connection error"))
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding team")
			},
		},
		{
			name: "given team exists but save fails when invoke then returns wrapped save error",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers"
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save error
				repository.On("Save", mock.Anything, mock.Anything).Return(fmt.Errorf("database error"))
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting match announcement in database")
			},
		},
		{
			name: "given empty team name when invoke then returns validation error",
			input: matchannouncement.Entity{
				TeamName:           "",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "team name cannot be empty")
			},
		},
		{
			name: "given day is in the past when invoke then returns validation error",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                yesterday,
				TimeSlot:           pastTimeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "day cannot be in the past")
			},
		},
		{
			name: "given incomplete location when invoke then returns validation error",
			input: matchannouncement.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           matchannouncement.NewLocation("", "", ""),
				AdmittedCategories: categoryRange,
				Status:             matchannouncement.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			given: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchannouncement.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "location must have country, province and locality")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			//set up
			repo := &mmocks.Repository{}
			teamRepo := &tmocks.Repository{}
			uc := usecases.NewCreateMatchAnnouncementUC(repo, teamRepo)

			// given
			tt.given(t, repo, teamRepo)

			// when
			result, err := uc.Invoke(context.Background(), tt.input)

			// then
			tt.then(t, result, err)
			repo.AssertExpectations(t)
			teamRepo.AssertExpectations(t)
		})
	}
}
