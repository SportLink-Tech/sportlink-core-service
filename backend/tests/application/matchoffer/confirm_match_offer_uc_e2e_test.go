package matchoffer_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	usecase "sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	dmatch "sportlink/api/domain/match"
	domain "sportlink/api/domain/matchoffer"
	dmatchrequest "sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/match"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"sportlink/tests/helper"

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
		then  func(t *testing.T, offerId string, entity *dmatch.Entity, err error)
	}{
		{
			name: "given a pending offer with pending requests when confirming then all the requests are rejected",
			setup: func(t *testing.T) usecase.ConfirmMatchOfferInput {
				ownerAcc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cabrerajjorge1@gmail.com").
					WithNickname("owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("jocabrera1@fi.uba.ar").
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
			then: func(t *testing.T, offerId string, entity *dmatch.Entity, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, entity)
				matchRequests, _ := mrRepo.Find(ctx, dmatchrequest.DomainQuery{
					MatchOfferIDs: []string{offerId},
					Statuses: []dmatchrequest.Status{
						dmatchrequest.StatusRejected,
					},
				})
				assert.Equal(t, entity.Status, dmatch.StatusAccepted)
				assert.True(t, len(matchRequests) == 1)
			},
		},
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
			then: func(t *testing.T, offerId string, entity *dmatch.Entity, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, entity.ID)
				assert.Equal(t, entity.Status, dmatch.StatusAccepted)
				page, _ := moRepo.Find(ctx, domain.DomainQuery{
					IDs: []string{offerId},
				})
				assert.Equal(t, domain.StatusConfirmed, page.Entities[0].Status)
			},
		},
		{
			name: "given a non-owner account when confirming an offer then it fails with unauthorized",
			setup: func(t *testing.T) usecase.ConfirmMatchOfferInput {
				ownerAcc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("confirm-unauth-owner@gmail.com").
					WithNickname("confirm-unauth-owner").
					Build(ctx)

				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithOwnerAccountID(ownerAcc.AccountID).
					WithCapacity(2).
					Build(ctx)

				return usecase.ConfirmMatchOfferInput{
					MatchOfferID:   offer.ID,
					OwnerAccountID: "wrong-account-id",
				}
			},
			then: func(t *testing.T, offerId string, entity *dmatch.Entity, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, entity)
			},
		},
		{
			name: "given an already confirmed offer when confirming again then it fails",
			setup: func(t *testing.T) usecase.ConfirmMatchOfferInput {
				ownerAcc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("confirm-double-owner@gmail.com").
					WithNickname("confirm-double-owner").
					Build(ctx)

				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithOwnerAccountID(ownerAcc.AccountID).
					WithCapacity(0).
					Build(ctx)

				input := usecase.ConfirmMatchOfferInput{
					MatchOfferID:   offer.ID,
					OwnerAccountID: ownerAcc.AccountID,
				}
				if _, err := confirmMatchOfferUC.Invoke(ctx, input); err != nil {
					t.Fatalf("failed to pre-confirm offer: %v", err)
				}

				return input
			},
			then: func(t *testing.T, offerId string, entity *dmatch.Entity, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, entity)
			},
		},
		{
			name: "given a pending offer with both accepted and pending requests when confirming then match includes only accepted participants",
			setup: func(t *testing.T) usecase.ConfirmMatchOfferInput {
				ownerAcc := helper.NewAccountBuilder(t, acRepo).
					WithEmail("confirm-mixed-owner@gmail.com").
					WithNickname("confirm-mixed-owner").
					Build(ctx)

				visitor1 := helper.NewAccountBuilder(t, acRepo).
					WithEmail("confirm-mixed-v1@fi.uba.ar").
					WithNickname("confirm-mixed-v1").
					Build(ctx)

				visitor2 := helper.NewAccountBuilder(t, acRepo).
					WithEmail("confirm-mixed-v2@fi.uba.ar").
					WithNickname("confirm-mixed-v2").
					Build(ctx)

				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithSport(common.Paddle).
					WithOwnerAccountID(ownerAcc.AccountID).
					WithCapacity(0).
					Build(ctx)

				// visitor1's request will be accepted
				request1 := helper.NewMatchRequestBuilder(t, mrRepo, moRepo).
					WithMatchOfferID(offer.ID).
					WithRequesterAccountID(visitor1.AccountID).
					Build(ctx)

				// visitor2's request stays pending
				_ = helper.NewMatchRequestBuilder(t, mrRepo, moRepo).
					WithMatchOfferID(offer.ID).
					WithRequesterAccountID(visitor2.AccountID).
					Build(ctx)

				// accept visitor1's request directly via repository
				if err := mrRepo.Save(ctx, request1.Accept()); err != nil {
					t.Fatalf("failed to accept request1: %v", err)
				}

				return usecase.ConfirmMatchOfferInput{
					MatchOfferID:   offer.ID,
					OwnerAccountID: ownerAcc.AccountID,
				}
			},
			then: func(t *testing.T, offerId string, entity *dmatch.Entity, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, entity)
				// owner + visitor1 (accepted) = 2; visitor2 (pending) is excluded
				assert.Len(t, entity.Participants, 2)

				// visitor2's pending request must be rejected after confirmation
				pendingRequests, _ := mrRepo.Find(ctx, dmatchrequest.DomainQuery{
					MatchOfferIDs: []string{offerId},
					Statuses:      []dmatchrequest.Status{dmatchrequest.StatusPending},
				})
				assert.Empty(t, pendingRequests)
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
			tc.then(t, input.MatchOfferID, entity, err)
		})
	}
}
