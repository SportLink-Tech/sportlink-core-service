package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/domain/team"
	mmocks "sportlink/mocks/api/domain/matchoffer"
	tmocks "sportlink/mocks/api/domain/team"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMatchOfferUC_Invoke(t *testing.T) {

	location := matchoffer.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	yesterday := time.Now().In(tz).AddDate(0, 0, -1)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 10, 0, 0, 0, tz)
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 12, 0, 0, 0, tz)

	timeSlot, _ := matchoffer.NewTimeSlot(startTime, endTime)
	pastTimeSlot, _ := matchoffer.NewTimeSlot(yesterdayStart, yesterdayEnd)
	categoryRange := matchoffer.NewSpecificCategories([]common.Category{5, 6, 7})
	greaterThanRange := matchoffer.NewGreaterThanCategory(5)

	tests := []struct {
		name  string
		input matchoffer.Entity
		on    func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository)
		then  func(t *testing.T, result *matchoffer.Entity, err error)
	}{
		{
			name: "save match offer successfully",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save offer
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity matchoffer.Entity) bool {
					return entity.TeamName == "Thunder Strikers" &&
						entity.Sport == common.Paddle &&
						entity.Status == matchoffer.StatusPending
				})).Return(nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "Thunder Strikers", result.TeamName)
				assert.Equal(t, common.Paddle, result.Sport)
			},
		},
		{
			name: "save match offer with GreaterThan category range successfully",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: greaterThanRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save offer
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity matchoffer.Entity) bool {
					return entity.TeamName == "Thunder Strikers" &&
						entity.AdmittedCategories.Type == matchoffer.RangeTypeGreaterThan
				})).Return(nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, matchoffer.RangeTypeGreaterThan, result.AdmittedCategories.Type)
			},
		},
		{
			name: "when team does not exist then offer cannot be created",
			input: matchoffer.Entity{
				TeamName:           "NonExistent Team",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team does not exist (empty slice)
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "NonExistent Team" &&
						reflect.DeepEqual(query.Sports, []common.Sport{common.Paddle})
				})).Return([]team.Entity{}, nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "team 'NonExistent Team' for sport 'Paddle' does not exist")
			},
		},
		{
			name: "when team repository throws an error then offer cannot be created",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team repository error
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers"
				})).Return([]team.Entity{}, fmt.Errorf("database connection error"))
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding team")
			},
		},
		{
			name: "when offer repository throws an error then it must be retrieved",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// Mock team exists
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return query.Name == "Thunder Strikers"
				})).Return([]team.Entity{{Name: "Thunder Strikers", Sport: common.Paddle}}, nil)

				// Mock save error
				repository.On("Save", mock.Anything, mock.Anything).Return(fmt.Errorf("database error"))
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting match offer in database")
			},
		},
		{
			name: "when team name is empty then offer cannot be created",
			input: matchoffer.Entity{
				TeamName:           "",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "team name cannot be empty")
			},
		},
		{
			name: "when day is in the past then offer cannot be created",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                yesterday,
				TimeSlot:           pastTimeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "day cannot be in the past")
			},
		},
		{
			name: "when location is empty then offer cannot be created",
			input: matchoffer.Entity{
				TeamName:           "Thunder Strikers",
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           matchoffer.NewLocation("", "", ""),
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository, teamRepository *tmocks.Repository) {
				// No mocks needed, validation happens before
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
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
			uc := usecases.NewCreateMatchOfferUC(repo, teamRepo)

			// given
			tt.on(t, repo, teamRepo)

			// when
			result, err := uc.Invoke(context.Background(), tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
