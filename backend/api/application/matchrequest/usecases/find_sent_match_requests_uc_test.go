package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchrequest"
	mrmocks "sportlink/mocks/api/domain/matchrequest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFindSentMatchRequestsUC(t *testing.T) {
	ctx := context.Background()
	requesterID := "requester-account"
	findErr := errors.New("repository find failed")

	tests := []struct {
		name     string
		statuses []matchrequest.Status
		given    func(t *testing.T, mrRepository *mrmocks.Repository)
		then     func(t *testing.T, result []matchrequest.Entity, err error)
	}{
		{
			name:     "given repository returns sent requests when invoke then returns entities",
			statuses: []matchrequest.Status{matchrequest.StatusPending},
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				entities := []matchrequest.Entity{
					{
						ID: "req-1", MatchAnnouncementID: "ann-1",
						OwnerAccountID: "owner-1", RequesterAccountID: requesterID,
						Status: matchrequest.StatusPending, CreatedAt: time.Now(),
					},
				}
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.RequesterAccountIDs) == 1 && q.RequesterAccountIDs[0] == requesterID &&
						len(q.Statuses) == 1 && q.Statuses[0] == matchrequest.StatusPending
				})).Return(entities, nil)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				assert.Equal(t, "req-1", result[0].ID)
				assert.Equal(t, requesterID, result[0].RequesterAccountID)
			},
		},
		{
			name:     "given repository returns empty slice when invoke then returns empty without error",
			statuses: nil,
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.RequesterAccountIDs) == 1 && q.RequesterAccountIDs[0] == requesterID && q.Statuses == nil
				})).Return([]matchrequest.Entity{}, nil)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name:     "given repository fails when invoke then returns wrapped error",
			statuses: []matchrequest.Status{matchrequest.StatusAccepted},
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.RequesterAccountIDs) == 1 && q.RequesterAccountIDs[0] == requesterID
				})).Return(nil, findErr)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding sent match requests")
				assert.ErrorIs(t, err, findErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mrRepository := &mrmocks.Repository{}
			uc := usecases.NewFindSentMatchRequestsUC(mrRepository)

			tt.given(t, mrRepository)

			result, err := uc.Invoke(ctx, requesterID, tt.statuses)

			tt.then(t, result, err)
			mrRepository.AssertExpectations(t)
		})
	}
}
