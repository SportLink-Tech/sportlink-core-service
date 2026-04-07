package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchrequest/usecases"
	domainoffer "sportlink/api/domain/matchoffer"
	domainreq "sportlink/api/domain/matchrequest"
	offermocks "sportlink/mocks/api/domain/matchoffer"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestCreateMatchRequestUC_Invoke(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name  string
		input usecases.CreateMatchRequestInput
		on    func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository)
		then  func(t *testing.T, result *domainreq.Entity, err error)
	}{
		{
			name: "given valid match offer when creating request then returns saved entity",
			input: usecases.CreateMatchRequestInput{
				MatchOfferID:       "offer-1",
				RequesterAccountID: "requester-acc",
			},
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				offer := domainoffer.Entity{
					ID:             "offer-1",
					OwnerAccountID: "owner-acc",
				}
				offerRepo.On("Find", mock.Anything, mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
					return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
				})).Return(domainoffer.Page{Entities: []domainoffer.Entity{offer}}, nil)
				reqRepo.On("Save", mock.Anything, mock.MatchedBy(func(e domainreq.Entity) bool {
					return e.MatchOfferID == "offer-1" &&
						e.OwnerAccountID == "owner-acc" &&
						e.RequesterAccountID == "requester-acc" &&
						e.Status == domainreq.StatusPending
				})).Return(nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "offer-1", result.MatchOfferID)
				assert.Equal(t, "owner-acc", result.OwnerAccountID)
				assert.Equal(t, "requester-acc", result.RequesterAccountID)
				assert.Equal(t, domainreq.StatusPending, result.Status)
				assert.Equal(t, domainreq.GenerateMatchRequestID("requester-acc", "offer-1"), result.ID)
			},
		},
		{
			name: "given match offer repository find fails when creating then returns wrapped error",
			input: usecases.CreateMatchRequestInput{
				MatchOfferID:       "offer-1",
				RequesterAccountID: "requester-acc",
			},
			on: func(t *testing.T, _ *reqmocks.Repository, offerRepo *offermocks.Repository) {
				offerRepo.On("Find", mock.Anything, mock.Anything).Return(domainoffer.Page{}, errors.New("db read error"))
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while finding match offer")
				assert.Contains(t, err.Error(), "db read error")
			},
		},
		{
			name: "given match offer not found when creating then returns not found error",
			input: usecases.CreateMatchRequestInput{
				MatchOfferID:       "missing-offer",
				RequesterAccountID: "requester-acc",
			},
			on: func(t *testing.T, _ *reqmocks.Repository, offerRepo *offermocks.Repository) {
				offerRepo.On("Find", mock.Anything, mock.Anything).Return(domainoffer.Page{Entities: []domainoffer.Entity{}}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match offer 'missing-offer' not found")
			},
		},
		{
			name: "given requester is offer owner when creating then returns error",
			input: usecases.CreateMatchRequestInput{
				MatchOfferID:       "offer-1",
				RequesterAccountID: "same-acc",
			},
			on: func(t *testing.T, _ *reqmocks.Repository, offerRepo *offermocks.Repository) {
				offer := domainoffer.Entity{ID: "offer-1", OwnerAccountID: "same-acc"}
				offerRepo.On("Find", mock.Anything, mock.Anything).Return(domainoffer.Page{Entities: []domainoffer.Entity{offer}}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "cannot send a match request to your own offer")
			},
		},
		{
			name: "given match request save fails when creating then returns wrapped error",
			input: usecases.CreateMatchRequestInput{
				MatchOfferID:       "offer-1",
				RequesterAccountID: "requester-acc",
			},
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				offer := domainoffer.Entity{ID: "offer-1", OwnerAccountID: "owner-acc"}
				offerRepo.On("Find", mock.Anything, mock.Anything).Return(domainoffer.Page{Entities: []domainoffer.Entity{offer}}, nil)
				reqRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("persist failed"))
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while saving match request")
				assert.Contains(t, err.Error(), "persist failed")
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			reqRepo := reqmocks.NewRepository(t)
			offerRepo := offermocks.NewRepository(t)
			uc := usecases.NewCreateMatchRequestUC(reqRepo, offerRepo)

			// given
			tt.on(t, reqRepo, offerRepo)

			// when
			result, err := uc.Invoke(ctx, tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
