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
	"sportlink/api/application/player/usecases"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/player"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/player"
	pmocks "sportlink/mocks/api/domain/player"
)

func TestPlayerCreationHandler(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name           string
		payloadRequest request.NewPlayerRequest
		on             func(t *testing.T, playerRepository *pmocks.Repository)
		assertions     func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "create a new player successfully",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Save", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.ID != "" && // ULID is generated
						entity.Category == common.L1 &&
						entity.Sport == common.Football
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"]) // ULID is generated
			},
		},
		{
			name: "create a new paddle player successfully",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Paddle",
				Category: 4,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"]) // ULID is generated
			},
		},
		{
			name: "create a new tennis player successfully",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Tennis",
				Category: 3,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"]) // ULID is generated
			},
		},
		{
			name: "create player with ULID generated automatically",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 2,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				// With ULID, each player gets a unique ID, so duplicates by ID are not possible
				// The test now just verifies successful creation
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"]) // ULID is generated
			},
		},
		{
			name: "fails when create a player with invalid category",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				Category: 99,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				// No repository calls expected for validation errors
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "invalid category value: 99")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name: "fails when create a player with invalid sport",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Basketball",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				// No repository calls expected for validation errors
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Error:Field validation for 'Sport' failed")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name: "create player without category defaults to Unranked",
			payloadRequest: request.NewPlayerRequest{
				Sport: "Football",
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Save", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.ID != "" && entity.Category == common.Unranked
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.NotEmpty(t, response["ID"]) // ULID is generated
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			playerRepository := new(pmocks.Repository)
			createPlayerUC := usecases.NewCreatePlayerUC(playerRepository)

			controller := player.NewController(&createPlayerUC, validator)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.POST("/player", controller.CreatePlayer)

			// given
			tc.on(t, playerRepository)
			jsonData, _ := json.Marshal(tc.payloadRequest)
			req, _ := http.NewRequest("POST", "/player", bytes.NewBuffer(jsonData))
			resp := httptest.NewRecorder()

			// when
			router.ServeHTTP(resp, req)

			// then
			response := createMapResponse(resp)
			tc.assertions(t, resp.Code, response)
		})
	}
}

func createMapResponse(resp *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	return response
}
