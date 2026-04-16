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
	"sportlink/tests/helper"
	"testing"
)

// cancelMatchRequest is a test helper that cancels a match request by directly
// updating the entity status in the repository, bypassing business-rule checks.
func cancelMatchRequest(t *testing.T, ctx context.Context, mrRepo dmatchrequest.Repository, entity *dmatchrequest.Entity) {
	t.Helper()
	cancelled := entity.Cancel()
	if err := mrRepo.Save(ctx, cancelled); err != nil {
		t.Fatalf("failed to cancel match request in setup: %v", err)
	}
}

func Test_FindMatchRequest(t *testing.T) {
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
	findMatchRequestUC := usecase.NewFindMatchRequestsUC(mrRepo)

	tests := []struct {
		name  string
		setup func(t *testing.T) *dmatchrequest.Entity
		query func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery
		then  func(t *testing.T, result []dmatchrequest.Entity, err error)
	}{
		{
			name: "given a match request created when find it then it must be retrieved",
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
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusPending, result[0].Status)
				assert.Regexp(t, `^AccountId#[^#]+#MatchOfferId#[^#]+$`, result[0].ID)
			},
		},
		{
			name: "given a match request created when find by requester account id then it must be retrieved",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-requester-owner@gmail.com").
					WithNickname("find-by-requester-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-requester-visitor@fi.uba.ar").
					WithNickname("find-by-requester-visitor").
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
			query: func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery {
				return dmatchrequest.DomainQuery{
					RequesterAccountIDs: []string{entity.RequesterAccountID},
				}
			},
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusPending, result[0].Status)
				assert.Regexp(t, `^AccountId#[^#]+#MatchOfferId#[^#]+$`, result[0].ID)
			},
		},
		{
			// Regression test: verifies that findByRequesterAccountID (used by "Solicitudes Enviadas")
			// returns requests with CANCEL status — the backend path behind
			// GET /account/:id/match-request?role=requester.
			name: "given a cancelled match request when find by requester account id then it must be retrieved with CANCEL status",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-cancelled-by-requester-owner@gmail.com").
					WithNickname("find-cancelled-by-requester-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-cancelled-by-requester-visitor@fi.uba.ar").
					WithNickname("find-cancelled-by-requester-visitor").
					Build(ctx)

				matchOffer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(2).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)

				entity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: visitor.AccountID,
				})
				cancelMatchRequest(t, ctx, mrRepo, entity)
				return entity
			},
			query: func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery {
				return dmatchrequest.DomainQuery{
					RequesterAccountIDs: []string{entity.RequesterAccountID},
				}
			},
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, dmatchrequest.StatusCancel, result[0].Status)
			},
		},
		{
			name: "given a match request created when find by match offer id then it must be retrieved",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-offer-owner@gmail.com").
					WithNickname("find-by-offer-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-offer-visitor@fi.uba.ar").
					WithNickname("find-by-offer-visitor").
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
			query: func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery {
				return dmatchrequest.DomainQuery{
					MatchOfferIDs: []string{entity.MatchOfferID},
				}
			},
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusPending, result[0].Status)
				assert.Regexp(t, `^AccountId#[^#]+#MatchOfferId#[^#]+$`, result[0].ID)
			},
		},
		{
			name: "given a match request created when find by owner account id then it must be retrieved",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-owner-owner@gmail.com").
					WithNickname("find-by-owner-owner").
					Build(ctx)

				visitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("find-by-owner-visitor@fi.uba.ar").
					WithNickname("find-by-owner-visitor").
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
			query: func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery {
				return dmatchrequest.DomainQuery{
					OwnerAccountIDs: []string{entity.OwnerAccountID},
				}
			},
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.True(t, len(result) == 1)
				assert.Equal(t, dmatchrequest.StatusPending, result[0].Status)
			},
		},
		{
			name: "given one pending and one cancelled request for the same owner when filtering by PENDING and ACCEPTED statuses then only the pending request is returned",
			setup: func(t *testing.T) *dmatchrequest.Entity {
				owner := helper.NewAccountBuilder(t, acRepo).
					WithEmail("status-filter-owner@gmail.com").
					WithNickname("status-filter-owner").
					Build(ctx)

				pendingVisitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("status-filter-pending@fi.uba.ar").
					WithNickname("status-filter-pending").
					Build(ctx)

				cancelledVisitor := helper.NewAccountBuilder(t, acRepo).
					WithEmail("status-filter-cancelled@fi.uba.ar").
					WithNickname("status-filter-cancelled").
					Build(ctx)

				matchOffer := helper.NewMatchOfferBuilder(t, moRepo).
					WithCapacity(5).
					WithOwnerAccountID(owner.AccountID).
					Build(ctx)

				// pending request — left as-is
				_, _ = createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: pendingVisitor.AccountID,
				})

				// cancelled request — created then manually cancelled
				cancelledEntity, _ := createMatchRequestUC.Invoke(ctx, usecase.CreateMatchRequestInput{
					MatchOfferID:       matchOffer.ID,
					RequesterAccountID: cancelledVisitor.AccountID,
				})
				cancelMatchRequest(t, ctx, mrRepo, cancelledEntity)

				// return a representative entity so the test harness can extract OwnerAccountID
				return cancelledEntity
			},
			query: func(entity *dmatchrequest.Entity) dmatchrequest.DomainQuery {
				return dmatchrequest.DomainQuery{
					OwnerAccountIDs: []string{entity.OwnerAccountID},
					Statuses:        []dmatchrequest.Status{dmatchrequest.StatusPending, dmatchrequest.StatusAccepted},
				}
			},
			then: func(t *testing.T, result []dmatchrequest.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, dmatchrequest.StatusPending, result[0].Status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchRequest := tt.setup(t)

			query := dmatchrequest.DomainQuery{IDs: []string{matchRequest.ID}}
			if tt.query != nil {
				query = tt.query(matchRequest)
			}

			// when
			result, err := findMatchRequestUC.Invoke(ctx, query)

			// then
			tt.then(t, result, err)
		})
	}
}
