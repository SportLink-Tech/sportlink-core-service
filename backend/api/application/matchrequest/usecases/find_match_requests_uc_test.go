package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchrequest/usecases"
	domainreq "sportlink/api/domain/matchrequest"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestFindMatchRequestsUC_Invoke(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name            string
		ownerAccountID  string
		on              func(t *testing.T, repository *reqmocks.Repository)
		then            func(t *testing.T, result []domainreq.Entity, err error)
	}{
		{
			name:           "given repository returns requests when finding by owner then returns entities",
			ownerAccountID: "owner-acc",
			on: func(t *testing.T, repository *reqmocks.Repository) {
				now := time.Now()
				entities := []domainreq.Entity{
					{
						ID:                 "mr-1",
						MatchOfferID:       "offer-1",
						OwnerAccountID:     "owner-acc",
						RequesterAccountID: "req-a",
						Status:             domainreq.StatusPending,
						CreatedAt:          now,
					},
				}
				repository.On("Find", mock.Anything, mock.MatchedBy(func(q domainreq.DomainQuery) bool {
					return len(q.OwnerAccountIDs) == 1 && q.OwnerAccountIDs[0] == "owner-acc"
				})).Return(entities, nil)
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, "mr-1", result[0].ID)
				assert.Equal(t, "owner-acc", result[0].OwnerAccountID)
			},
		},
		{
			name:           "given repository returns empty when finding by owner then returns empty slice",
			ownerAccountID: "owner-without-requests",
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(q domainreq.DomainQuery) bool {
					return len(q.OwnerAccountIDs) == 1 && q.OwnerAccountIDs[0] == "owner-without-requests"
				})).Return([]domainreq.Entity{}, nil)
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:           "given repository find fails when finding by owner then returns wrapped error",
			ownerAccountID: "owner-acc",
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("Find", mock.Anything, mock.Anything).Return(nil, errors.New("query failed"))
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding match requests")
				assert.Contains(t, err.Error(), "query failed")
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			repository := reqmocks.NewRepository(t)
			uc := usecases.NewFindMatchRequestsUC(repository)

			// given
			tt.on(t, repository)

			// when
			result, err := uc.Invoke(ctx, tt.ownerAccountID)

			// then
			tt.then(t, result, err)
		})
	}
}
