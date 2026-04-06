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
		name  string
		query domainreq.DomainQuery
		on    func(t *testing.T, repository *reqmocks.Repository)
		then  func(t *testing.T, result []domainreq.Entity, err error)
	}{
		{
			name:  "given owner account id when finding then returns received requests",
			query: domainreq.DomainQuery{OwnerAccountIDs: []string{"owner-acc"}},
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
				assert.Equal(t, "owner-acc", result[0].OwnerAccountID)
			},
		},
		{
			name:  "given requester account id and status when finding then returns sent requests",
			query: domainreq.DomainQuery{RequesterAccountIDs: []string{"requester-acc"}, Statuses: []domainreq.Status{domainreq.StatusPending}},
			on: func(t *testing.T, repository *reqmocks.Repository) {
				now := time.Now()
				entities := []domainreq.Entity{
					{
						ID:                 "mr-2",
						MatchOfferID:       "offer-2",
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
			name:  "given repository find fails then returns wrapped error",
			query: domainreq.DomainQuery{OwnerAccountIDs: []string{"owner-acc"}},
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
			repository := reqmocks.NewRepository(t)
			uc := usecases.NewFindMatchRequestsUC(repository)

			tt.on(t, repository)

			result, err := uc.Invoke(ctx, tt.query)

			tt.then(t, result, err)
		})
	}
}
