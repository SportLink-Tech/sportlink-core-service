package matchrequest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domain "sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/middleware"
	cmr "sportlink/api/infrastructure/rest/matchrequest"
	mocks "sportlink/mocks/api/infrastructure/rest/matchrequest"
)

func TestFindSentMatchRequests(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name    string
		acc     string
		rawPath string
		given   func(t *testing.T, create *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase)
		then    func(t *testing.T, code int, body interface{})
	}{
		{
			name:    "given valid statuses query when get sent then invokes use case with parsed statuses",
			acc:     "req-acc",
			rawPath: "/account/req-acc/sent-match-request?statuses=PENDING,ACCEPTED",
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				now := time.Now()
				sent.On("Invoke", mock.Anything, "req-acc", []domain.Status{domain.StatusPending, domain.StatusAccepted}).Return([]domain.Entity{
					{ID: "s1", RequesterAccountID: "req-acc", Status: domain.StatusPending, CreatedAt: now},
				}, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
			},
		},
		{
			name:    "given invalid status in query when get sent then returns bad request",
			acc:     "req-acc",
			rawPath: "/account/req-acc/sent-match-request?statuses=NOT_A_STATUS",
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
			},
		},
		{
			name:    "given use case fails when get sent then returns conflict",
			acc:     "req-acc",
			rawPath: "/account/req-acc/sent-match-request",
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				sent.On("Invoke", mock.Anything, "req-acc", mock.MatchedBy(func(s []domain.Status) bool {
					return s == nil || len(s) == 0
				})).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusConflict, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", m["code"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			createM := &mocks.CreateMatchRequestUseCase{}
			findM := &mocks.FindMatchRequestsUseCase{}
			sentM := &mocks.FindSentMatchRequestsUseCase{}
			updM := &mocks.UpdateMatchRequestStatusUseCase{}
			tc.given(t, createM, findM, sentM, updM)

			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			ctl := cmr.NewController(createM, findM, sentM, updM, v)
			r.GET("/account/:accountId/sent-match-request", ctl.FindSentMatchRequests)

			req := httptest.NewRequest(http.MethodGet, tc.rawPath, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &body)
			tc.then(t, rec.Code, body)

			createM.AssertExpectations(t)
			findM.AssertExpectations(t)
			sentM.AssertExpectations(t)
			updM.AssertExpectations(t)
		})
	}
}
