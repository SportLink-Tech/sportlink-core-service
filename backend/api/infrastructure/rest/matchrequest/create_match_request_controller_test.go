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

	"sportlink/api/application/matchrequest/usecases"
	domain "sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/middleware"
	cmr "sportlink/api/infrastructure/rest/matchrequest"
	mocks "sportlink/mocks/api/infrastructure/rest/matchrequest"
)

func TestCreateMatchRequest(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name  string
		accID string
		annID string
		given func(t *testing.T, create *mocks.CreateMatchRequestUseCase, find *mocks.FindMatchRequestsUseCase, sent *mocks.FindSentMatchRequestsUseCase, upd *mocks.UpdateMatchRequestStatusUseCase)
		then  func(t *testing.T, code int, body map[string]interface{})
	}{
		{
			name:  "given use case succeeds when posting match request then returns created",
			accID: "requester-acc",
			annID: "announcement-01",
			given: func(t *testing.T, create *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				now := time.Now()
				ent := &domain.Entity{
					ID: "mr-1", MatchAnnouncementID: "announcement-01",
					OwnerAccountID: "owner-acc", RequesterAccountID: "requester-acc",
					Status: domain.StatusPending, CreatedAt: now,
				}
				create.On("Invoke", mock.Anything, mock.MatchedBy(func(in usecases.CreateMatchRequestInput) bool {
					return in.MatchAnnouncementID == "announcement-01" && in.RequesterAccountID == "requester-acc"
				})).Return(ent, nil)
			},
			then: func(t *testing.T, code int, body map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, code)
				assert.Equal(t, "mr-1", body["id"])
				assert.Equal(t, "announcement-01", body["match_announcement_id"])
				assert.Equal(t, "owner-acc", body["owner_account_id"])
				assert.Equal(t, "requester-acc", body["requester_account_id"])
				assert.Equal(t, "PENDING", body["status"])
			},
		},
		{
			name:  "given use case fails when posting match request then returns conflict",
			accID: "requester-acc",
			annID: "announcement-01",
			given: func(t *testing.T, create *mocks.CreateMatchRequestUseCase, _ *mocks.FindMatchRequestsUseCase, _ *mocks.FindSentMatchRequestsUseCase, _ *mocks.UpdateMatchRequestStatusUseCase) {
				create.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
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
			r.POST("/account/:accountId/match-announcement/:announcementId/match-request", ctl.CreateMatchRequest)

			req := httptest.NewRequest(http.MethodPost, "/account/"+tc.accID+"/match-announcement/"+tc.annID+"/match-request", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body map[string]interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &body)
			tc.then(t, rec.Code, body)

			createM.AssertExpectations(t)
			findM.AssertExpectations(t)
			sentM.AssertExpectations(t)
			updM.AssertExpectations(t)
		})
	}
}
