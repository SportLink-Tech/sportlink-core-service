package matchoffer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	dmatchoffer "sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/dev/testcontainer"
	"sportlink/tests/helper"
)

func Test_FindMatchOfferUC(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)

	moRepo := matchoffer.NewRepository(dynamoDbClient, "SportLinkCore")

	tests := []struct {
		name  string
		given func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery
		then  func(t *testing.T, result *usecases.FindMatchOfferResult, err error)
	}{
		{
			name: "given paddle and football offers when querying by paddle then returns only paddle offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Football).Build(ctx)
				return dmatchoffer.DomainQuery{Sports: []common.Sport{common.Paddle}}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, result.Page.Total)
				for _, e := range result.Entities {
					assert.Equal(t, common.Paddle, e.Sport)
				}
			},
		},
		{
			name: "given offers from different owners when querying by owner then returns only that owner's offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithOwnerAccountID("owner-a").Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithOwnerAccountID("owner-a").Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithOwnerAccountID("owner-b").Build(ctx)
				return dmatchoffer.DomainQuery{OwnerAccountID: "owner-a"}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, result.Page.Total)
				for _, e := range result.Entities {
					assert.Equal(t, "owner-a", e.OwnerAccountID)
				}
			},
		},
		{
			name: "given offers with different sports when querying by multiple sports then returns all matching offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Tennis).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Football).Build(ctx)
				return dmatchoffer.DomainQuery{Sports: []common.Sport{common.Paddle, common.Tennis}}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, result.Page.Total)
			},
		},
		{
			name: "given pending and confirmed offers when querying by confirmed status then returns only confirmed offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithStatus(dmatchoffer.StatusPending).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithStatus(dmatchoffer.StatusPending).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithStatus(dmatchoffer.StatusConfirmed).Build(ctx)
				return dmatchoffer.DomainQuery{Statuses: []dmatchoffer.Status{dmatchoffer.StatusConfirmed}}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, result.Page.Total)
				assert.Equal(t, dmatchoffer.StatusConfirmed, result.Entities[0].Status)
			},
		},
		{
			name: "given offers on different days when querying by date range then returns only offers within range",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				tz := dmatchoffer.NewLocation("Argentina", "Buenos Aires", "Palermo").GetTimezone()
				tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
				nextWeek := time.Now().In(tz).AddDate(0, 0, 7)
				nextWeekStart := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 18, 0, 0, 0, tz)
				nextWeekEnd := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 20, 0, 0, 0, tz)

				helper.NewMatchOfferBuilder(t, moRepo).
					WithDay(tomorrow).
					WithStatus(dmatchoffer.StatusConfirmed).
					Build(ctx) // tomorrow by default

				helper.NewMatchOfferBuilder(t, moRepo).
					WithDay(nextWeek).
					WithStatus(dmatchoffer.StatusCancelled).
					WithTimeSlot(nextWeekStart, nextWeekEnd).
					Build(ctx)

				inFiveDays := time.Now().In(tz).AddDate(0, 0, 5)
				return dmatchoffer.DomainQuery{
					FromDate: tomorrow.Add(-time.Hour),
					ToDate:   inFiveDays,
				}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, result.Page.Total)
				assert.Equal(t, dmatchoffer.StatusConfirmed, result.Entities[0].Status)
			},
		},
		{
			name: "given offers in different localities when querying by locality then returns only matching offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithLocation("Argentina", "Buenos Aires", "Palermo").Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithLocation("Argentina", "Buenos Aires", "Palermo").Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithLocation("Argentina", "Santa Fe", "Rosario").Build(ctx)
				loc := dmatchoffer.NewLocation("Argentina", "Buenos Aires", "Palermo")
				return dmatchoffer.DomainQuery{Location: &loc}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, result.Page.Total)
				for _, e := range result.Entities {
					assert.Equal(t, "Palermo", e.Location.Locality)
				}
			},
		},
		{
			name: "given three offers when querying with limit two then returns two entities with correct page info",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				return dmatchoffer.DomainQuery{Sports: []common.Sport{common.Paddle}, Limit: 2}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Len(t, result.Entities, 2)
				assert.Equal(t, 3, result.Page.Total)
				assert.Equal(t, 2, result.Page.OutOf)
				assert.Equal(t, 1, result.Page.Number)
			},
		},
		{
			name: "given three offers when querying second page with limit two and offset two then returns remaining offer",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Paddle).Build(ctx)
				return dmatchoffer.DomainQuery{Sports: []common.Sport{common.Paddle}, Limit: 2, Offset: 2}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Len(t, result.Entities, 1)
				assert.Equal(t, 3, result.Page.Total)
				assert.Equal(t, 2, result.Page.Number)
			},
		},
		{
			name: "given offers saved with SaveAll when querying by owner and sport then returns all saved offers",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				ownerAccountId := "owner-account-id"
				singlePaddle := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithCapacity(2).
					WithOwnerAccountID(ownerAccountId).
					WithCategoryRange(dmatchoffer.CategoryRange{
						Categories: []common.Category{
							common.L1, common.L3, common.L5,
						},
					}).
					Build(ctx)

				commonPaddle := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithCapacity(4).
					WithOwnerAccountID(ownerAccountId).
					WithCategoryRange(dmatchoffer.CategoryRange{
						Categories: []common.Category{
							common.L1, common.L3, common.L5,
						},
					}).
					Build(ctx)

				moRepo.SaveAll(ctx, []dmatchoffer.Entity{*singlePaddle, *commonPaddle})

				return dmatchoffer.DomainQuery{
					OwnerAccountID: ownerAccountId,
					Sports:         []common.Sport{common.Paddle},
				}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(result.Entities))
			},
		},
		{
			name: "given no matching offers when querying by sport then returns empty result",
			given: func(t *testing.T, moRepo dmatchoffer.Repository) dmatchoffer.DomainQuery {
				helper.NewMatchOfferBuilder(t, moRepo).WithSport(common.Football).Build(ctx)
				return dmatchoffer.DomainQuery{Sports: []common.Sport{common.Tennis}}
			},
			then: func(t *testing.T, result *usecases.FindMatchOfferResult, err error) {
				assert.Nil(t, err)
				assert.Empty(t, result.Entities)
				assert.Equal(t, 0, result.Page.Total)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			t.Cleanup(func() {
				testcontainer.ClearDynamoDbTable(t, dynamoDbClient, "SportLinkCore")
			})
			matchOfferUC := usecases.NewFindMatchOfferUC(moRepo)

			// given
			query := tc.given(t, moRepo)

			// when
			result, err := matchOfferUC.Invoke(ctx, query)

			// then
			tc.then(t, result, err)
		})
	}
}
