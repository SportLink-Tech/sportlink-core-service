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

func TestFindSentMatchRequestsUC_Invoke(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name               string
		requesterAccountID string
		statuses           []domainreq.Status
		on                 func(t *testing.T, repository *reqmocks.Repository)
		then               func(t *testing.T, result []domainreq.Entity, err error)
	}{
		{
			name:               "given repository returns sent requests when finding then returns entities",
			requesterAccountID: "requester-acc",
			statuses:           []domainreq.Status{domainreq.StatusPending},
			on: func(t *testing.T, repository *reqmocks.Repository) {
				now := time.Now()
				entities := []domainreq.Entity{
					{
						ID:                 "mr-1",
						MatchOfferID:       "offer-1",
						OwnerAccountID:     "owner-acc",
						RequesterAccountID: "requester-acc",
						Status:             domainreq.StatusPending,
						CreatedAt:          now,
					},
				}
				repository.On("Find", mock.Anything, mock.MatchedBy(func(q domainreq.DomainQuery) bool {
					return len(q.RequesterAccountIDs) == 1 &&
						q.RequesterAccountIDs[0] == "requester-acc" &&
						len(q.Statuses) == 1 &&
						q.Statuses[0] == domainreq.StatusPending
				})).Return(entities, nil)
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, "requester-acc", result[0].RequesterAccountID)
			},
		},
		{
			name:               "given nil statuses when finding sent then passes nil statuses to repository",
			requesterAccountID: "requester-acc",
			statuses:           nil,
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(q domainreq.DomainQuery) bool {
					return len(q.RequesterAccountIDs) == 1 &&
						q.RequesterAccountIDs[0] == "requester-acc" &&
						q.Statuses == nil
				})).Return([]domainreq.Entity{}, nil)
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:               "given repository find fails when finding sent then returns wrapped error",
			requesterAccountID: "requester-acc",
			statuses:           nil,
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("Find", mock.Anything, mock.Anything).Return(nil, errors.New("scan error"))
			},
			then: func(t *testing.T, result []domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding sent match requests")
				assert.Contains(t, err.Error(), "scan error")
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			repository := reqmocks.NewRepository(t)
			uc := usecases.NewFindSentMatchRequestsUC(repository)

			// given
			tt.on(t, repository)

			// when
			result, err := uc.Invoke(ctx, tt.requesterAccountID, tt.statuses)

			// then
			tt.then(t, result, err)
		})
	}
}
