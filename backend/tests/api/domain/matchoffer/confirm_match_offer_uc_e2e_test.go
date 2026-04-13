package matchoffer_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	usecase "sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	dmatch "sportlink/api/domain/match"
	domain "sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/match"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"sportlink/tests/api/helper"

	"testing"
)

func Test_ConfirmMatchOfferUC(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)

	// repositories
	mRepo := match.NewRepository(dynamoDbClient, "SportLinkCore")
	moRepo := matchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	mrRepo := matchrequest.NewRepository(dynamoDbClient, "SportLinkCore")
	acRepo := account.NewRepository(dynamoDbClient, "SportLinkCore")

	confirmMatchOfferUC := usecase.NewConfirmMatchOfferUC(mRepo, moRepo, mrRepo)

	tests := []struct {
		name  string
		setup func(t *testing.T) usecase.ConfirmMatchOfferInput
		then  func(t *testing.T, entity *dmatch.Entity, err error)
	}{
		{
			name: "given a pending offer with accepted requests when confirming then creates match",
			setup: func(t *testing.T) usecase.ConfirmMatchOfferInput {
				ownerAcc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cabrerajjorge@gmail.com").
					WithNickname("owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("jocabrera@fi.uba.ar").
					WithNickname("visitor").
					Build(ctx)

				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithOwnerAccountID(ownerAcc.AccountID).
					WithCapacity(2).
					Build(ctx)

				_ = helper.NewMatchRequestBuilder(t, mrRepo, moRepo).
					WithMatchOfferID(offer.ID).
					WithRequesterAccountID(visitor.AccountID).
					Build(ctx)

				return usecase.ConfirmMatchOfferInput{
					MatchOfferID:   offer.ID,
					OwnerAccountID: ownerAcc.AccountID,
				}
			},
			then: func(t *testing.T, entity *dmatch.Entity, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, entity.ID)
				assert.Equal(t, entity.Status, dmatch.StatusAccepted)
				page, _ := moRepo.Find(ctx, domain.DomainQuery{
					IDs: []string{entity.ID},
				})
				assert.Equal(t, domain.StatusConfirmed, page.Entities[0].Status)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			input := tc.setup(t)

			// when
			entity, err := confirmMatchOfferUC.Invoke(ctx, input)

			// then
			tc.then(t, entity, err)
		})
	}
}
