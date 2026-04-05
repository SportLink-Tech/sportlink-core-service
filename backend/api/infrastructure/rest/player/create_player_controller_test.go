package player_test

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

	"sportlink/api/application/player/request"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/player"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/player"
	amocks "sportlink/mocks/api/application"
)

func TestCreatePlayer(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name           string
		payloadRequest request.NewPlayerRequest
		given          func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity])
		then           func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "given use case succeeds when creating player then returns created",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 1,
			},
			given: func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity]) {
				out := domain.NewPlayer(common.L1, common.Football)
				uc.On("Invoke", mock.Anything, mock.MatchedBy(func(e domain.Entity) bool {
					return e.Category == common.L1 && e.Sport == common.Football
				})).Return(&out, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"])
			},
		},
		{
			name: "given use case fails when creating player then returns conflict",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 1,
			},
			given: func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity]) {
				uc.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
		{
			name: "given invalid category when creating then returns bad request",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 99,
			},
			given: func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity]) {},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Contains(t, response["message"], "invalid category value")
			},
		},
		{
			name: "given invalid sport when creating then returns bad request",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Basketball",
				Category: 1,
			},
			given: func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity]) {},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Contains(t, response["message"], "Sport")
			},
		},
		{
			name: "given empty category defaults when creating then invokes use case with unranked",
			payloadRequest: request.NewPlayerRequest{
				Sport: "Football",
			},
			given: func(t *testing.T, uc *amocks.UseCase[domain.Entity, domain.Entity]) {
				out := domain.NewPlayer(common.Unranked, common.Football)
				uc.On("Invoke", mock.Anything, mock.MatchedBy(func(e domain.Entity) bool {
					return e.Category == common.Unranked && e.Sport == common.Football
				})).Return(&out, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := amocks.NewUseCase[domain.Entity, domain.Entity](t)
			tc.given(t, uc)

			controller := player.NewController(uc, v)
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.POST("/player", controller.CreatePlayer)

			jsonData, _ := json.Marshal(tc.payloadRequest)
			req, _ := http.NewRequest(http.MethodPost, "/player", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var response map[string]interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &response)
			tc.then(t, rec.Code, response)
		})
	}
}
