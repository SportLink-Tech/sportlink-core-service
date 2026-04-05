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

func TestFindMatchRequests(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name  string
		owner string
		given func(t *testing.T, create *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase)
		then  func(t *testing.T, code int, body interface{})
	}{
		{
			name:  "given use case returns requests when get then returns ok json array",
			owner: "owner-acc",
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				now := time.Now()
				find.On("Invoke", mock.Anything, "owner-acc").Return([]domain.Entity{
					{ID: "r1", MatchAnnouncementID: "a1", OwnerAccountID: "owner-acc", RequesterAccountID: "x", Status: domain.StatusPending, CreatedAt: now},
				}, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
				first := arr[0].(map[string]interface{})
				assert.Equal(t, "r1", first["id"])
				assert.Equal(t, "owner-acc", first["owner_account_id"])
			},
		},
		{
			name:  "given use case fails when get then returns conflict",
			owner: "owner-acc",
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				find.On("Invoke", mock.Anything, "owner-acc").Return(nil, assert.AnError)
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
			r.GET("/account/:accountId/match-request", ctl.FindMatchRequests)

			req := httptest.NewRequest(http.MethodGet, "/account/"+tc.owner+"/match-request", nil)
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
