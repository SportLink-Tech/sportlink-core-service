package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/domain/matchrequest"
	mamocks "sportlink/mocks/api/domain/matchannouncement"
	mrmocks "sportlink/mocks/api/domain/matchrequest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreateMatchRequestUC(t *testing.T) {
	ctx := context.Background()

	const (
		announcementID     = "01HZXK8QTESTANNOUNCE01"
		ownerAccountID     = "owner-account"
		requesterAccountID = "requester-account"
	)

	findRepoErr := errors.New("dynamodb query failed")
	saveRepoErr := errors.New("conditional check failed")

	tests := []struct {
		name  string
		input usecases.CreateMatchRequestInput
		given func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository)
		then  func(t *testing.T, result *matchrequest.Entity, err error)
	}{
		{
			name: "given a match announcement exists and requester is not the owner when invoke then returns pending request and saves",
			input: usecases.CreateMatchRequestInput{
				MatchAnnouncementID: announcementID,
				RequesterAccountID:  requesterAccountID,
			},
			given: func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository) {
				ann := matchannouncement.Entity{
					ID:             announcementID,
					OwnerAccountID: ownerAccountID,
				}
				maRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchannouncement.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == announcementID
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{ann},
					Total:    1,
				}, nil)

				mrRepository.On("Save", mock.Anything, mock.MatchedBy(func(e matchrequest.Entity) bool {
					return e.MatchAnnouncementID == announcementID &&
						e.OwnerAccountID == ownerAccountID &&
						e.RequesterAccountID == requesterAccountID &&
						e.Status == matchrequest.StatusPending &&
						e.ID != ""
				})).Return(nil)
			},
			then: func(t *testing.T, result *matchrequest.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, announcementID, result.MatchAnnouncementID)
				assert.Equal(t, ownerAccountID, result.OwnerAccountID)
				assert.Equal(t, requesterAccountID, result.RequesterAccountID)
				assert.Equal(t, matchrequest.StatusPending, result.Status)
				assert.NotEmpty(t, result.ID)
				assert.False(t, result.CreatedAt.IsZero())
			},
		},
		{
			name: "given the match announcement repository fails when invoke then returns wrapped find error and does not save",
			input: usecases.CreateMatchRequestInput{
				MatchAnnouncementID: announcementID,
				RequesterAccountID:  requesterAccountID,
			},
			given: func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository) {
				maRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchannouncement.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == announcementID
				})).Return(matchannouncement.Page{}, findRepoErr)
			},
			then: func(t *testing.T, result *matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding match announcement")
				assert.ErrorIs(t, err, findRepoErr)
			},
		},
		{
			name: "given no match announcement is found when invoke then returns not found error and does not save",
			input: usecases.CreateMatchRequestInput{
				MatchAnnouncementID: announcementID,
				RequesterAccountID:  requesterAccountID,
			},
			given: func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository) {
				maRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchannouncement.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == announcementID
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{},
					Total:    0,
				}, nil)
			},
			then: func(t *testing.T, result *matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "not found")
				assert.Contains(t, err.Error(), announcementID)
			},
		},
		{
			name: "given the requester is the announcement owner when invoke then returns error and does not save",
			input: usecases.CreateMatchRequestInput{
				MatchAnnouncementID: announcementID,
				RequesterAccountID:  ownerAccountID,
			},
			given: func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository) {
				ann := matchannouncement.Entity{
					ID:             announcementID,
					OwnerAccountID: ownerAccountID,
				}
				maRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchannouncement.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == announcementID
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{ann},
					Total:    1,
				}, nil)
			},
			then: func(t *testing.T, result *matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "cannot send a match request to your own announcement")
			},
		},
		{
			name: "given saving the match request fails when invoke then returns wrapped save error",
			input: usecases.CreateMatchRequestInput{
				MatchAnnouncementID: announcementID,
				RequesterAccountID:  requesterAccountID,
			},
			given: func(t *testing.T, maRepository *mamocks.Repository, mrRepository *mrmocks.Repository) {
				ann := matchannouncement.Entity{
					ID:             announcementID,
					OwnerAccountID: ownerAccountID,
				}
				maRepository.On("Find", mock.Anything, mock.MatchedBy(func(q matchannouncement.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == announcementID
				})).Return(matchannouncement.Page{
					Entities: []matchannouncement.Entity{ann},
					Total:    1,
				}, nil)

				mrRepository.On("Save", mock.Anything, mock.MatchedBy(func(e matchrequest.Entity) bool {
					return e.MatchAnnouncementID == announcementID &&
						e.OwnerAccountID == ownerAccountID &&
						e.RequesterAccountID == requesterAccountID &&
						e.Status == matchrequest.StatusPending
				})).Return(saveRepoErr)
			},
			then: func(t *testing.T, result *matchrequest.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while saving match request")
				assert.ErrorIs(t, err, saveRepoErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			maRepository := &mamocks.Repository{}
			mrRepository := &mrmocks.Repository{}
			uc := usecases.NewCreateMatchRequestUC(mrRepository, maRepository)

			tt.given(t, maRepository, mrRepository)

			result, err := uc.Invoke(ctx, tt.input)

			tt.then(t, result, err)
			maRepository.AssertExpectations(t)
			mrRepository.AssertExpectations(t)
		})
	}
}
