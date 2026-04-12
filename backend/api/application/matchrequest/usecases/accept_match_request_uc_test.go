package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	matchofferevent "sportlink/api/application/matchoffer/events"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/common"
	domainoffer "sportlink/api/domain/matchoffer"
	domainreq "sportlink/api/domain/matchrequest"
	eventmocks "sportlink/mocks/api/application/events"
	offermocks "sportlink/mocks/api/domain/matchoffer"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestAcceptMatchRequestUC_Invoke(t *testing.T) {
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

	validInput := usecases.AcceptMatchRequestInput{
		MatchRequestId: pendingRequest.ID,
		OwnerAccountID: "owner-1",
	}

	testCases := []struct {
		name  string
		input usecases.AcceptMatchRequestInput
		on    func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent])
		then  func(t *testing.T, result *domainreq.Entity, err error)
	}{
		{
			name:  "given pending request and pending offer with no capacity when accepting then saves accepted request",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
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
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domainreq.StatusAccepted, result.Status)
				assert.Equal(t, pendingRequest.ID, result.ID)
			},
		},
		{
			name: "given non-owner account when accepting then returns unauthorized",
			input: usecases.AcceptMatchRequestInput{
				MatchRequestId: pendingRequest.ID,
				OwnerAccountID: "another-account",
			},
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "owner account ID does not match")
			},
		},
		{
			name: "given match request not found when accepting then returns error",
			input: usecases.AcceptMatchRequestInput{
				MatchRequestId: "non-existent-id",
				OwnerAccountID: "owner-1",
			},
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "non-existent-id"
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
			name:  "given repository error when finding match request then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
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
			name:  "given already accepted match request when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
				acceptedRequest := pendingRequest
				acceptedRequest.Status = domainreq.StatusAccepted
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match request is not pending")
			},
		},
		{
			name:  "given match offer not found when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
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
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{}, Total: 0}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match offer not found")
			},
		},
		{
			name:  "given already confirmed match offer when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
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
				assert.Contains(t, err.Error(), "match offer is not pending")
			},
		},
		{
			name:  "given repository fails when saving accepted request then returns error",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
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
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(errors.New("request save failed"))
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "request save failed")
			},
		},
		{
			name:  "given offer with capacity reached when accepting last request then publishes event",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
				offerWithCapacity := pendingOffer
				offerWithCapacity.Capacity = 2 // owner + 1 requester

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
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{offerWithCapacity}, Total: 1}, nil)

				reqRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(nil)

				// After save, count accepted requests (now 1 accepted == capacity-1)
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1" &&
							len(q.Statuses) == 1 && q.Statuses[0] == domainreq.StatusAccepted
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				publisher.On("Publish",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(e matchofferevent.MatchOfferCapacityReachedEvent) bool {
						return e.MatchOfferID == "offer-1" && e.OwnerAccountID == "owner-1"
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domainreq.StatusAccepted, result.Status)
			},
		},
		{
			name:  "given offer with capacity not yet reached when accepting request then does not publish event",
			input: validInput,
			on: func(t *testing.T, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository, publisher *eventmocks.Publisher[matchofferevent.MatchOfferCapacityReachedEvent]) {
				offerWithCapacity := pendingOffer
				offerWithCapacity.Capacity = 4 // owner + 3 requesters

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
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{offerWithCapacity}, Total: 1}, nil)

				reqRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(nil)

				// count returns 1 accepted, capacity-1 == 3, not reached
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1" &&
							len(q.Statuses) == 1 && q.Statuses[0] == domainreq.StatusAccepted
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)
			},
			then: func(t *testing.T, result *domainreq.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, domainreq.StatusAccepted, result.Status)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reqRepo := reqmocks.NewRepository(t)
			offerRepo := offermocks.NewRepository(t)
			publisher := eventmocks.NewPublisher[matchofferevent.MatchOfferCapacityReachedEvent](t)
			uc := usecases.NewAcceptMatchRequestUC(reqRepo, offerRepo, publisher)

			tc.on(t, reqRepo, offerRepo, publisher)

			result, err := uc.Invoke(ctx, tc.input)

			tc.then(t, result, err)
		})
	}
}
