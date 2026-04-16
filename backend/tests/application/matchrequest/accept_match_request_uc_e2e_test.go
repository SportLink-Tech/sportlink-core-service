package matchrequest_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	matchofferevent "sportlink/api/application/matchoffer/events"
	usecase "sportlink/api/application/matchrequest/usecases"
	dmatchrequest "sportlink/api/domain/matchrequest"
	ievents "sportlink/api/infrastructure/events"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"sportlink/tests/helper"
	"testing"
)

func Test_AcceptMatchRequest(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)

	// repositories
	moRepo := matchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	mrRepo := matchrequest.NewRepository(dynamoDbClient, "SportLinkCore")
	acRepo := account.NewRepository(dynamoDbClient, "SportLinkCore")
	capacityPublisher := ievents.NewChannelPublisher[matchofferevent.MatchOfferCapacityReachedEvent](100)

	// use cases
	createMatchRequestUC := usecase.NewCreateMatchRequestUC(mrRepo, moRepo)
	acceptMatchRequestUC := usecase.NewAcceptMatchRequestUC(mrRepo, moRepo, capacityPublisher)
	findMatchRequestUC := usecase.NewFindMatchRequestsUC(mrRepo)

	tests := []struct {
		name        string
		setup       func(t *testing.T) *dmatchrequest.Entity
		acceptInput func(entity *dmatchrequest.Entity) usecase.AcceptMatchRequestInput
		then        func(t *testing.T, acceptErr error, result []dmatchrequest.Entity, err error)
	}{
		{
			name: "given a match request accepted when find it by status accepted then it must be retrieved",
			setup: func(t *testing.T) *dmatchrequest.Entity {
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

				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				return entity
			},
			then: func(t *testing.T, acceptErr error, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, acceptErr)
				assert.Nil(t, err)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusAccepted, result[0].Status)
				assert.Regexp(t, `^AccountId#[^#]+#MatchOfferId#[^#]+$`, result[0].ID)
			},
		},
		{
			name: "given a non-owner account when accepting a match request then it fails with unauthorized",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("accept-unauth-owner@gmail.com").
					WithNickname("accept-unauth-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("accept-unauth-visitor@fi.uba.ar").
					WithNickname("accept-unauth-visitor").
					Build(ctx)

				matchOffer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(2).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)

				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				return entity
			},
			acceptInput: func(entity *dmatchrequest.Entity) usecase.AcceptMatchRequestInput {
				return usecase.AcceptMatchRequestInput{
					MatchRequestId: entity.ID,
					OwnerAccountID: "wrong-account-id",
				}
			},
			then: func(t *testing.T, acceptErr error, result []dmatchrequest.Entity, err error) {
				assert.NotNil(t, acceptErr)
				assert.Nil(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name: "given an already accepted match request when accepting again then it fails",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("accept-double-owner@gmail.com").
					WithNickname("accept-double-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("accept-double-visitor@fi.uba.ar").
					WithNickname("accept-double-visitor").
					Build(ctx)

				matchOffer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(0).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)

				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: visitor.AccountID,
				})

				// pre-accept the request
				accepted, _ := acceptMatchRequestUC.Invoke(ctx, usecase.AcceptMatchRequestInput{
					MatchRequestId: entity.ID,
					OwnerAccountID: entity.OwnerAccountID,
				})
				return accepted
			},
			then: func(t *testing.T, acceptErr error, result []dmatchrequest.Entity, err error) {
				assert.NotNil(t, acceptErr)
				assert.Nil(t, err)
				// the pre-accepted request is still there
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusAccepted, result[0].Status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchRequest := tt.setup(t)

			input := usecase.AcceptMatchRequestInput{
				MatchRequestId: matchRequest.ID,
				OwnerAccountID: matchRequest.OwnerAccountID,
			}
			if tt.acceptInput != nil {
				input = tt.acceptInput(matchRequest)
			}

			// given
			_, acceptErr := acceptMatchRequestUC.Invoke(ctx, input)

			// when
			result, err := findMatchRequestUC.Invoke(ctx, dmatchrequest.DomainQuery{
				RequesterAccountIDs: []string{matchRequest.RequesterAccountID},
				Statuses:            []dmatchrequest.Status{dmatchrequest.StatusAccepted},
			})

			// then
			tt.then(t, acceptErr, result, err)
		})
	}
}
