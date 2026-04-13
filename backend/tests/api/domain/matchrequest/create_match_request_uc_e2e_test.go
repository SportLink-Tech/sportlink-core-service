package matchrequest__test

import (
	"context"
	"github.com/stretchr/testify/assert"
	dmatchrequest "sportlink/api/domain/matchrequest"
	"sportlink/tests/api/helper"

	usecase "sportlink/api/application/matchrequest/usecases"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"testing"
)

func Test_CreateMatchRequest(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)

	// repositories
	moRepo := matchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	mrRepo := matchrequest.NewRepository(dynamoDbClient, "SportLinkCore")
	acRepo := account.NewRepository(dynamoDbClient, "SportLinkCore")

	uc := usecase.NewCreateMatchRequestUC(mrRepo, moRepo)

	tests := []struct {
		name  string
		setup func(t *testing.T) usecase.CreateMatchRequestInput
		then  func(t *testing.T, result *dmatchrequest.Entity, err error)
	}{
		{
			name: "given an account and a pending match offer when create a match request then the request is created",
			setup: func(t *testing.T) usecase.CreateMatchRequestInput {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cabrerajjorge@gmail.com").
					WithNickname("owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("jocabrera@fi.uba.ar").
					WithNickname("visitor").
					Build(ctx)

				matchOffer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(2).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)

				return usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: visitor.AccountID,
				}
			},
			then: func(t *testing.T, result *dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, result.ID)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.setup(t)

			// when
			result, err := uc.Invoke(ctx, input)

			// then
			tt.then(t, result, err)
		})
	}
}
