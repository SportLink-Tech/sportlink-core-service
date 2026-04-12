package usecases_test

import (
	"context"
	"fmt"
	"sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/matchoffer"
	mmocks "sportlink/mocks/api/domain/matchoffer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMatchOfferUC_Invoke(t *testing.T) {
	ctx := context.Background()

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
		on    func(t *testing.T, repository *mmocks.Repository)
		then  func(t *testing.T, result *matchoffer.Entity, err error)
	}{
		{
			name: "given valid offer when saving then returns saved entity",
			input: matchoffer.Entity{
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(entity matchoffer.Entity) bool {
						return entity.Sport == common.Paddle && entity.Status == matchoffer.StatusPending
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, common.Paddle, result.Sport)
			},
		},
		{
			name: "given offer with team name when saving then saves successfully",
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
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(entity matchoffer.Entity) bool {
						return entity.TeamName == "Thunder Strikers"
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "Thunder Strikers", result.TeamName)
			},
		},
		{
			name: "given offer with GreaterThan category range when saving then saves successfully",
			input: matchoffer.Entity{
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: greaterThanRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(entity matchoffer.Entity) bool {
						return entity.AdmittedCategories.Type == matchoffer.RangeTypeGreaterThan
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, matchoffer.RangeTypeGreaterThan, result.AdmittedCategories.Type)
			},
		},
		{
			name: "given repository error when saving then returns wrapped error",
			input: matchoffer.Entity{
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.Anything,
				).Return(fmt.Errorf("database error"))
			},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting match offer in database")
			},
		},
		{
			name: "given day in the past when creating then returns error",
			input: matchoffer.Entity{
				Sport:              common.Paddle,
				Day:                yesterday,
				TimeSlot:           pastTimeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "day cannot be in the past")
			},
		},
		{
			name: "given empty location when creating then returns error",
			input: matchoffer.Entity{
				Sport:              common.Paddle,
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           matchoffer.NewLocation("", "", ""),
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "location must have country, province and locality")
			},
		},
		{
			name: "given empty sport when creating then returns error",
			input: matchoffer.Entity{
				Sport:              "",
				Day:                tomorrow,
				TimeSlot:           timeSlot,
				Location:           location,
				AdmittedCategories: categoryRange,
				Status:             matchoffer.StatusPending,
				CreatedAt:          time.Now().In(tz),
			},
			on: func(t *testing.T, repository *mmocks.Repository) {},
			then: func(t *testing.T, result *matchoffer.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "sport cannot be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mmocks.NewRepository(t)
			uc := usecases.NewCreateMatchOfferUC(repo)

			tt.on(t, repo)

			result, err := uc.Invoke(ctx, tt.input)

			tt.then(t, result, err)
		})
	}
}
