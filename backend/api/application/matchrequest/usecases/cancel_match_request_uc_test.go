package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/common"
	domainoffer "sportlink/api/domain/matchoffer"
	domainreq "sportlink/api/domain/matchrequest"
	offermocks "sportlink/mocks/api/domain/matchoffer"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestCancelMatchRequestUC_Invoke(t *testing.T) {
	ctx := context.Background()

	fixedDay := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	fixedNow := time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC)

	pendingRequest := domainreq.Entity{
		ID:                 "AccountId#requester-1#MatchOfferId#offer-1",
		MatchOfferID:       "offer-1",
		OwnerAccountID:     "owner-1",
		RequesterAccountID: "requester-1",
		Status:             domainreq.StatusPending,
		CreatedAt:          fixedNow,
	}

	pendingOffer := domainoffer.Entity{
		ID:       "offer-1",
		TeamName: "Los Leones FC",
		Sport:    common.Paddle,
		Day:      fixedDay,
		TimeSlot: domainoffer.TimeSlot{
			StartTime: fixedDay.Add(18 * time.Hour),
			EndTime:   fixedDay.Add(20 * time.Hour),
		},
		Status:         domainoffer.StatusPending,
		OwnerAccountID: "owner-1",
	}

	validInput := usecases.CancelMatchRequestInput{
		MatchRequestId:     pendingRequest.ID,
		RequesterAccountID: "requester-1",
	}

	testCases := []struct {
		name  string
		input usecases.CancelMatchRequestInput
		on    func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository)
		then  func(t *testing.T, result *domainreq.Entity, err error)
	}{
		{
			name:  "given pending request and pending offer when cancelling then saves cancelled request",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusCancel
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domainreq.StatusCancel, result.Status)
				assert.Equal(t, pendingRequest.ID, result.ID)
			},
		},
		{
			name:  "given match request not found when cancelling then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match request not found")
			},
		},
		{
			name:  "given repository error when finding request then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return(nil, errors.New("db connection error"))
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "db connection error")
			},
		},
		{
			name:  "given already rejected request when cancelling then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				rejectedRequest := pendingRequest
				rejectedRequest.Status = domainreq.StatusRejected
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{rejectedRequest}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match request is already rejected")
			},
		},
		{
			name:  "given confirmed offer when cancelling request then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				confirmedOffer := pendingOffer
				confirmedOffer.Status = domainoffer.StatusConfirmed
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{confirmedOffer}, Total: 1}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match offer is already confirmed")
			},
		},
		{
			name:  "given repository fails when saving cancelled request then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusCancel
					}),
				).Return(errors.New("save failed"))
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "save failed")
			},
		},
		{
			name:  "given accepted request when cancelling then saves cancelled request",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				acceptedRequest := pendingRequest
				acceptedRequest.Status = domainreq.StatusAccepted
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)

				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusCancel
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domainreq.StatusCancel, result.Status)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reqRepo := reqmocks.NewRepository(t)
			offerRepo := offermocks.NewRepository(t)
			uc := usecases.NewCancelMatchRequestUC(reqRepo, offerRepo)

			tc.on(t, reqRepo, offerRepo)

			result, err := uc.Invoke(ctx, tc.input)

			tc.then(t, result, err)
		})
	}
}
