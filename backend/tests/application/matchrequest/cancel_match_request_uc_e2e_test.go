package matchrequest_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	usecase "sportlink/api/application/matchrequest/usecases"
	dmatchrequest "sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/api/infrastructure/persistence/matchoffer"
	"sportlink/api/infrastructure/persistence/matchrequest"
	"sportlink/dev/testcontainer"
	"sportlink/tests/api/helper"
	"testing"
)

func Test_CancelMatchRequest(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)

	// repositories
	moRepo := matchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	mrRepo := matchrequest.NewRepository(dynamoDbClient, "SportLinkCore")
	acRepo := account.NewRepository(dynamoDbClient, "SportLinkCore")

	// use cases
	createMatchRequestUC := usecase.NewCreateMatchRequestUC(mrRepo, moRepo)
	cancelMatchRequestUC := usecase.NewCancelMatchRequestUC(mrRepo, moRepo)
	findMatchRequestUC := usecase.NewFindMatchRequestsUC(mrRepo)

	tests := []struct {
		name  string
		setup func(t *testing.T) *dmatchrequest.Entity
		then  func(t *testing.T, cancelErr error, result []dmatchrequest.Entity)
	}{
		{
			name: "given a pending match request when cancelling then status is CANCEL",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-owner@gmail.com").
					WithNickname("cancel-owner").
					Build(ctx)
				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-visitor@fi.uba.ar").
					WithNickname("cancel-visitor").
					Build(ctx)
				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(0).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)
				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       offer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				return entity
			},
			then: func(t *testing.T, cancelErr error, result []dmatchrequest.Entity) {
				assert.Nil(t, cancelErr)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusCancel, result[0].Status)
				assert.Regexp(t, `^AccountId#[^#]+#MatchOfferId#[^#]+$`, result[0].ID)
			},
		},
		{
			name: "given a rejected match request when cancelling then it fails",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-rejected-owner@gmail.com").
					WithNickname("cancel-rejected-owner").
					Build(ctx)
				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-rejected-visitor@fi.uba.ar").
					WithNickname("cancel-rejected-visitor").
					Build(ctx)
				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(0).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)
				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       offer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				if err := mrRepo.Save(ctx, entity.Reject()); err != nil {
					t.Fatalf("failed to reject request: %v", err)
				}
				return entity
			},
			then: func(t *testing.T, cancelErr error, result []dmatchrequest.Entity) {
				assert.NotNil(t, cancelErr)
				assert.Contains(t, cancelErr.Error(), "match request is already rejected")
			},
		},
		{
			name: "given a confirmed match offer when cancelling the request then it fails",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-confirmed-owner@gmail.com").
					WithNickname("cancel-confirmed-owner").
					Build(ctx)
				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("cancel-confirmed-visitor@fi.uba.ar").
					WithNickname("cancel-confirmed-visitor").
					Build(ctx)
				offer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(0).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)
				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       offer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				if err := moRepo.Save(ctx, offer.Confirm()); err != nil {
					t.Fatalf("failed to confirm offer: %v", err)
				}
				return entity
			},
			then: func(t *testing.T, cancelErr error, result []dmatchrequest.Entity) {
				assert.NotNil(t, cancelErr)
				assert.Contains(t, cancelErr.Error(), "match offer is already confirmed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchRequest := tt.setup(t)

			_, cancelErr := cancelMatchRequestUC.Invoke(ctx, usecase.CancelMatchRequestInput{
				MatchRequestId:     matchRequest.ID,
				RequesterAccountID: matchRequest.RequesterAccountID,
			})

			result, _ := findMatchRequestUC.Invoke(ctx, dmatchrequest.DomainQuery{
				IDs: []string{matchRequest.ID},
			})

			tt.then(t, cancelErr, result)
		})
	}
}
