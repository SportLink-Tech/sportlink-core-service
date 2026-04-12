package matchoffer_test

import (
	"context"
	usecase "sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	dmatch "sportlink/api/domain/match"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/match"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"sportlink/tests/api/helper"

	"testing"
)

func TestConfirmMatchOfferUC(t *testing.T) {
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
				acc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cabrerajjorge@gmail.com").
					WithNickname("testuser").
					Build(ctx)

				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithOwnerAccountID(acc.AccountID).
					WithCapacity(2).
					Build(ctx)

				return usecase.ConfirmMatchOfferInput{
					MatchOfferID:   offer.ID,
					OwnerAccountID: acc.AccountID,
				}
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
