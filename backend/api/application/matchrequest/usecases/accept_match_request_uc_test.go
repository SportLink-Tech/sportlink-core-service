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
	domainmatch "sportlink/api/domain/match"
	domainoffer "sportlink/api/domain/matchoffer"
	domainreq "sportlink/api/domain/matchrequest"
	matchmocks "sportlink/mocks/api/domain/match"
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
		on    func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository)
		then  func(t *testing.T, result *domainmatch.Entity, err error)
	}{
		{
			name:  "given pending request and pending offer when accepting then creates match and updates request and offer",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				matchRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return m.LocalAccountID == "owner-1" &&
							m.VisitorAccountID == "requester-1" &&
							m.Sport == common.Paddle &&
							m.Day.Equal(fixedDay) &&
							m.Status == domainmatch.StatusAccepted
					}),
				).Return(nil)

				reqRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(nil)

				offerRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(o domainoffer.Entity) bool {
						return o.ID == "offer-1" && o.Status == domainoffer.StatusConfirmed
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "owner-1", result.LocalAccountID)
				assert.Equal(t, "requester-1", result.VisitorAccountID)
				assert.Equal(t, common.Paddle, result.Sport)
				assert.Equal(t, domainmatch.StatusAccepted, result.Status)
			},
		},
		{
			name: "given non-owner account when accepting then returns unauthorized",
			input: usecases.AcceptMatchRequestInput{
				MatchRequestId: pendingRequest.ID,
				OwnerAccountID: "another-account",
			},
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
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
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "non-existent-id"
					}),
				).Return([]domainreq.Entity{}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match request not found")
			},
		},
		{
			name:  "given repository error when finding match request then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return(nil, errors.New("db connection error"))
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "db connection error")
			},
		},
		{
			name:  "given already accepted match request when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				acceptedRequest := pendingRequest
				acceptedRequest.Status = domainreq.StatusAccepted
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match request is not pending")
			},
		},
		{
			name:  "given match offer not found when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{}, Total: 0}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match offer not found")
			},
		},
		{
			name:  "given already confirmed match offer when accepting then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				confirmedOffer := pendingOffer
				confirmedOffer.Status = domainoffer.StatusConfirmed
				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{confirmedOffer}, Total: 1}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match offer is not pending")
			},
		},
		{
			name:  "given match repository fails when saving match then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				matchRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return m.LocalAccountID == "owner-1" && m.VisitorAccountID == "requester-1"
					}),
				).Return(errors.New("match save failed"))
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "match save failed")
			},
		},
		{
			name:  "given match request repository fails when saving accepted request then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				matchRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return m.LocalAccountID == "owner-1" && m.VisitorAccountID == "requester-1"
					}),
				).Return(nil)

				reqRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(errors.New("request save failed"))
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "request save failed")
			},
		},
		{
			name:  "given match offer repository fails when saving confirmed offer then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, reqRepo *reqmocks.Repository, offerRepo *offermocks.Repository) {
				reqRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == pendingRequest.ID
					}),
				).Return([]domainreq.Entity{pendingRequest}, nil)

				offerRepo.On("Find",
					mock.Anything,
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				matchRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return m.LocalAccountID == "owner-1" && m.VisitorAccountID == "requester-1"
					}),
				).Return(nil)

				reqRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(r domainreq.Entity) bool {
						return r.ID == pendingRequest.ID && r.Status == domainreq.StatusAccepted
					}),
				).Return(nil)

				offerRepo.On("Save",
					mock.Anything,
					mock.MatchedBy(func(o domainoffer.Entity) bool {
						return o.ID == "offer-1" && o.Status == domainoffer.StatusConfirmed
					}),
				).Return(errors.New("offer save failed"))
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "offer save failed")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matchRepo := matchmocks.NewRepository(t)
			reqRepo := reqmocks.NewRepository(t)
			offerRepo := offermocks.NewRepository(t)
			uc := usecases.NewAcceptMatchRequestUC(matchRepo, reqRepo, offerRepo)

			tc.on(t, matchRepo, reqRepo, offerRepo)

			result, err := uc.Invoke(ctx, tc.input)

			tc.then(t, result, err)
		})
	}
}
