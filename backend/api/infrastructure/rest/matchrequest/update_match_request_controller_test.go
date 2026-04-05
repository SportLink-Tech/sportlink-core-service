package matchrequest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchrequest/usecases"
	domain "sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/middleware"
	cmr "sportlink/api/infrastructure/rest/matchrequest"
	mocks "sportlink/mocks/api/infrastructure/rest/matchrequest"
)

func TestUpdateMatchRequestStatus(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name    string
		acc     string
		reqID   string
		payload map[string]string
		rawJSON string
		given   func(t *testing.T, create *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase)
		then    func(t *testing.T, code int, body map[string]interface{})
	}{
		{
			name:    "given valid body when patch then returns no content",
			acc:     "owner-acc",
			reqID:   "mr-99",
			payload: map[string]string{"status": "ACCEPTED"},
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase) {
				upd.On("Invoke", mock.Anything, mock.MatchedBy(func(in usecases.UpdateMatchRequestStatusInput) bool {
					return in.ID == "mr-99" && in.OwnerAccountID == "owner-acc" && in.NewStatus == domain.StatusAccepted
				})).Return(nil)
			},
			then: func(t *testing.T, code int, body map[string]interface{}) {
				assert.Equal(t, http.StatusNoContent, code)
				assert.Nil(t, body)
			},
		},
		{
			name:    "given invalid json when patch then returns bad request",
			acc:     "owner-acc",
			reqID:   "mr-99",
			rawJSON: `{not-json`,
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
			},
			then: func(t *testing.T, code int, body map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				assert.Equal(t, "invalid_request_format", body["code"])
			},
		},
		{
			name:    "given invalid status value when patch then returns validation error",
			acc:     "owner-acc",
			reqID:   "mr-99",
			payload: map[string]string{"status": "PENDING"},
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
			},
			then: func(t *testing.T, code int, body map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				assert.Equal(t, "request_validation_failed", body["code"])
			},
		},
		{
			name:    "given use case fails when patch then returns conflict",
			acc:     "owner-acc",
			reqID:   "mr-99",
			payload: map[string]string{"status": "REJECTED"},
			given: func(t *testing.T, _ *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase) {
				upd.On("Invoke", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			then: func(t *testing.T, code int, body map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, code)
				assert.Equal(t, "use_case_execution_error", body["code"])
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
			r.PATCH("/account/:accountId/match-request/:requestId", ctl.UpdateMatchRequestStatus)

			var bodyBytes []byte
			if tc.rawJSON != "" {
				bodyBytes = []byte(tc.rawJSON)
			} else {
				bodyBytes, _ = json.Marshal(tc.payload)
			}
			req := httptest.NewRequest(http.MethodPatch, "/account/"+tc.acc+"/match-request/"+tc.reqID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body map[string]interface{}
			if rec.Body.Len() > 0 {
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
			}
			tc.then(t, rec.Code, body)

			createM.AssertExpectations(t)
			findM.AssertExpectations(t)
			sentM.AssertExpectations(t)
			updM.AssertExpectations(t)
		})
	}
}
