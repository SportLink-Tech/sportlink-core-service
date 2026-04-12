package match_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/match/usecases"
	"sportlink/api/domain/common"
	domainmatch "sportlink/api/domain/match"
	"sportlink/api/infrastructure/middleware"
	cmatches "sportlink/api/infrastructure/rest/match"
	amocks "sportlink/mocks/api/application"
)

type FindMatchesUCMock = amocks.UseCase[usecases.FindMatchesInput, []domainmatch.Entity]

func TestFindMatches(t *testing.T) {
	fixedDay := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)

	acceptedMatch := domainmatch.Entity{
		ID:           "01MATCH001",
		Participants: []string{"owner-1", "requester-1"},
		Sport:        common.Paddle,
		Day:          fixedDay,
		Status:       domainmatch.StatusAccepted,
		CreatedAt:    fixedDay,
	}

	testCases := []struct {
		name        string
		accountID   string
		queryParams map[string]string
		on          func(t *testing.T, uc *FindMatchesUCMock)
		then        func(t *testing.T, code int, body interface{})
	}{
		{
			name:        "given no status filter when listing matches then returns all account matches",
			accountID:   "owner-1",
			queryParams: map[string]string{},
			on: func(t *testing.T, uc *FindMatchesUCMock) {
				matches := []domainmatch.Entity{acceptedMatch}
				uc.On("Invoke",
					mock.Anything,
					mock.MatchedBy(func(in usecases.FindMatchesInput) bool {
						return in.AccountID == "owner-1" && len(in.Statuses) == 0
					}),
				).Return(&matches, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
				first := arr[0].(map[string]interface{})
				assert.Equal(t, "01MATCH001", first["id"])
				participants := first["participants"].([]interface{})
				assert.Equal(t, "owner-1", participants[0])
				assert.Equal(t, "Paddle", first["sport"])
				assert.Equal(t, "ACCEPTED", first["status"])
			},
		},
		{
			name:        "given single status filter when listing matches then returns filtered matches",
			accountID:   "owner-1",
			queryParams: map[string]string{"statuses": "ACCEPTED"},
			on: func(t *testing.T, uc *FindMatchesUCMock) {
				matches := []domainmatch.Entity{acceptedMatch}
				uc.On("Invoke",
					mock.Anything,
					mock.MatchedBy(func(in usecases.FindMatchesInput) bool {
						return in.AccountID == "owner-1" &&
							len(in.Statuses) == 1 &&
							in.Statuses[0] == domainmatch.StatusAccepted
					}),
				).Return(&matches, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
			},
		},
		{
			name:        "given multiple status filter when listing matches then returns filtered matches",
			accountID:   "owner-1",
			queryParams: map[string]string{"statuses": "ACCEPTED,PLAYED"},
			on: func(t *testing.T, uc *FindMatchesUCMock) {
				matches := []domainmatch.Entity{acceptedMatch}
				uc.On("Invoke",
					mock.Anything,
					mock.MatchedBy(func(in usecases.FindMatchesInput) bool {
						return in.AccountID == "owner-1" &&
							len(in.Statuses) == 2 &&
							in.Statuses[0] == domainmatch.StatusAccepted &&
							in.Statuses[1] == domainmatch.StatusPlayed
					}),
				).Return(&matches, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
			},
		},
		{
			name:        "given no matches when listing then returns empty array",
			accountID:   "owner-1",
			queryParams: map[string]string{},
			on: func(t *testing.T, uc *FindMatchesUCMock) {
				empty := []domainmatch.Entity{}
				uc.On("Invoke",
					mock.Anything,
					mock.MatchedBy(func(in usecases.FindMatchesInput) bool {
						return in.AccountID == "owner-1"
					}),
				).Return(&empty, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 0)
			},
		},
		{
			name:        "given invalid status when listing matches then returns validation error",
			accountID:   "owner-1",
			queryParams: map[string]string{"statuses": "INVALID"},
			on:          func(t *testing.T, uc *FindMatchesUCMock) {},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
				assert.Contains(t, m["message"], "invalid status")
			},
		},
		{
			name:        "given repository error when listing matches then returns conflict",
			accountID:   "owner-1",
			queryParams: map[string]string{},
			on: func(t *testing.T, uc *FindMatchesUCMock) {
				uc.On("Invoke",
					mock.Anything,
					mock.MatchedBy(func(in usecases.FindMatchesInput) bool {
						return in.AccountID == "owner-1"
					}),
				).Return(nil, errors.New("db connection error"))
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusConflict, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", m["code"])
				assert.Contains(t, m["message"], "db connection error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ucMock := amocks.NewUseCase[usecases.FindMatchesInput, []domainmatch.Entity](t)
			controller := cmatches.NewController(ucMock)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(middleware.ErrorHandler())
			r.GET("/account/:account_id/match", controller.FindMatches)

			tc.on(t, ucMock)

			req := httptest.NewRequest(http.MethodGet, "/account/"+tc.accountID+"/match", nil)
			q := req.URL.Query()
			for k, v := range tc.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			var body interface{}
			if rec.Body.Len() > 0 {
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
			}

			tc.then(t, rec.Code, body)
		})
	}
}
