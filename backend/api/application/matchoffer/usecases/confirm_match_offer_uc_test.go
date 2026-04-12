package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	domainmatch "sportlink/api/domain/match"
	domainoffer "sportlink/api/domain/matchoffer"
	domainreq "sportlink/api/domain/matchrequest"
	matchmocks "sportlink/mocks/api/domain/match"
	offermocks "sportlink/mocks/api/domain/matchoffer"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestConfirmMatchOfferUC_Invoke(t *testing.T) {
	ctx := context.Background()

	fixedDay := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	fixedNow := time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC)

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

	acceptedRequest := domainreq.Entity{
		ID:                 "AccountId#requester-1#MatchOfferId#offer-1",
		MatchOfferID:       "offer-1",
		OwnerAccountID:     "owner-1",
		RequesterAccountID: "requester-1",
		Status:             domainreq.StatusAccepted,
		CreatedAt:          fixedNow,
	}

	validInput := usecases.ConfirmMatchOfferInput{
		MatchOfferID:   "offer-1",
		OwnerAccountID: "owner-1",
	}

	testCases := []struct {
		name  string
		input usecases.ConfirmMatchOfferInput
		on    func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository)
		then  func(t *testing.T, result *domainmatch.Entity, err error)
	}{
		{
			name:  "given pending offer with accepted requests when confirming then creates match and confirms offer",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1" &&
							len(q.Statuses) == 1 && q.Statuses[0] == domainreq.StatusAccepted
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)

				matchRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return len(m.Participants) == 2 &&
							m.Participants[0] == "owner-1" &&
							m.Participants[1] == "requester-1" &&
							m.Sport == common.Paddle &&
							m.Day.Equal(fixedDay) &&
							m.Status == domainmatch.StatusAccepted
					}),
				).Return(nil)

				offerRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(o domainoffer.Entity) bool {
						return o.ID == "offer-1" && o.Status == domainoffer.StatusConfirmed
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, []string{"owner-1", "requester-1"}, result.Participants)
				assert.Equal(t, common.Paddle, result.Sport)
				assert.Equal(t, domainmatch.StatusAccepted, result.Status)
			},
		},
		{
			name:  "given pending offer with multiple accepted requests when confirming then creates match with all participants",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				twoRequests := []domainreq.Entity{
					acceptedRequest,
					{
						ID: "AccountId#requester-2#MatchOfferId#offer-1", MatchOfferID: "offer-1",
						OwnerAccountID: "owner-1", RequesterAccountID: "requester-2",
						Status: domainreq.StatusAccepted, CreatedAt: fixedNow,
					},
				}
				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1"
					}),
				).Return(twoRequests, nil)

				matchRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return len(m.Participants) == 3 &&
							m.Participants[0] == "owner-1" &&
							m.Participants[1] == "requester-1" &&
							m.Participants[2] == "requester-2"
					}),
				).Return(nil)

				offerRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(o domainoffer.Entity) bool {
						return o.ID == "offer-1" && o.Status == domainoffer.StatusConfirmed
					}),
				).Return(nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, []string{"owner-1", "requester-1", "requester-2"}, result.Participants)
			},
		},
		{
			name: "given non-owner account when confirming then returns unauthorized",
			input: usecases.ConfirmMatchOfferInput{
				MatchOfferID:   "offer-1",
				OwnerAccountID: "another-account",
			},
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "owner account ID does not match")
			},
		},
		{
			name:  "given offer not found when confirming then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
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
			name:  "given repository error when finding offer then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{}, errors.New("db connection error"))
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "db connection error")
			},
		},
		{
			name:  "given already confirmed offer when confirming then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				confirmedOffer := pendingOffer
				confirmedOffer.Status = domainoffer.StatusConfirmed
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
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
			name:  "given no accepted requests when confirming then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1"
					}),
				).Return([]domainreq.Entity{}, nil)
			},
			then: func(t *testing.T, result *domainmatch.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "no accepted requests found")
			},
		},
		{
			name:  "given repository error when finding accepted requests then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1"
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
			name:  "given match repository fails when saving match then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1"
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)

				matchRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return len(m.Participants) == 2
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
			name:  "given offer repository fails when saving confirmed offer then returns error",
			input: validInput,
			on: func(t *testing.T, matchRepo *matchmocks.Repository, offerRepo *offermocks.Repository, reqRepo *reqmocks.Repository) {
				offerRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainoffer.DomainQuery) bool {
						return len(q.IDs) == 1 && q.IDs[0] == "offer-1"
					}),
				).Return(domainoffer.Page{Entities: []domainoffer.Entity{pendingOffer}, Total: 1}, nil)

				reqRepo.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q domainreq.DomainQuery) bool {
						return len(q.MatchOfferIDs) == 1 && q.MatchOfferIDs[0] == "offer-1"
					}),
				).Return([]domainreq.Entity{acceptedRequest}, nil)

				matchRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(m domainmatch.Entity) bool {
						return len(m.Participants) == 2
					}),
				).Return(nil)

				offerRepo.On("Save",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
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
			offerRepo := offermocks.NewRepository(t)
			reqRepo := reqmocks.NewRepository(t)
			uc := usecases.NewConfirmMatchOfferUC(matchRepo, offerRepo, reqRepo)

			tc.on(t, matchRepo, offerRepo, reqRepo)

			result, err := uc.Invoke(ctx, tc.input)

			tc.then(t, result, err)
		})
	}
}
