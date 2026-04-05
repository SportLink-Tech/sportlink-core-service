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

func TestNewFindMatchRequestsUC(t *testing.T) {
	ctx := context.Background()
	ownerID := "owner-account"
	findErr := errors.New("repository find failed")

	tests := []struct {
		name  string
		given func(t *testing.T, mrRepository *mrmocks.Repository)
		then  func(t *testing.T, result []matchrequest.Entity, err error)
	}{
		{
			name: "given repository returns incoming requests when invoke then returns entities",
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				entities := []matchrequest.Entity{
					{
						ID: "req-1", MatchAnnouncementID: "ann-1",
						OwnerAccountID: ownerID, RequesterAccountID: "other",
						Status: matchrequest.StatusPending, CreatedAt: time.Now(),
					},
				}
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.OwnerAccountIDs) == 1 && q.OwnerAccountIDs[0] == ownerID
				})).Return(entities, nil)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				assert.Equal(t, ownerID, result[0].OwnerAccountID)
			},
		},
		{
			name: "given repository returns no requests when invoke then returns empty without error",
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.OwnerAccountIDs) == 1 && q.OwnerAccountIDs[0] == ownerID
				})).Return([]matchrequest.Entity{}, nil)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name: "given repository fails when invoke then returns wrapped error",
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchrequest.DomainQuery) bool {
					return len(q.OwnerAccountIDs) == 1 && q.OwnerAccountIDs[0] == ownerID
				})).Return(nil, findErr)
			},
			then: func(t *testing.T, result []matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding match requests")
				assert.ErrorIs(t, err, findErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mrRepository := &mrmocks.Repository{}
			uc := usecases.NewFindMatchRequestsUC(mrRepository)

			tt.given(t, mrRepository)

			result, err := uc.Invoke(ctx, ownerID)

			tt.then(t, result, err)
			mrRepository.AssertExpectations(t)
		})
	}
}
