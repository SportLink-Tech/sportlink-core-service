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
				ID:       "player1",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return query.Id == "player1" && query.Category == common.L1 && query.Sport == common.Football
				})).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.ID == "player1" &&
						entity.Category == common.L1 &&
						entity.Sport == common.Football
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "player1", response["ID"])
			},
		},
		{
			name: "create a new paddle player successfully",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Paddle",
				ID:       "paddle_player1",
				Category: 4,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "paddle_player1", response["ID"])
			},
		},
		{
			name: "create a new tennis player successfully",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Tennis",
				ID:       "tennis_player1",
				Category: 3,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name: "fails when player already exists",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				ID:       "existing_player",
				Category: 2,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				existingPlayer := domain.Entity{
					ID:       "existing_player",
					Category: common.L2,
					Sport:    common.Football,
				}
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{existingPlayer}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
				assert.Contains(t, response["message"], "Player already exist")
			},
		},
		{
			name: "fails when create a player with invalid category",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				ID:       "player2",
				Category: 99,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
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
				ID:       "player3",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Error:Field validation for 'Sport' failed")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name: "fails when ID is missing",
			payloadRequest: request.NewPlayerRequest{
				Sport:    "Football",
				ID:       "",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Error:Field validation for 'ID'")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name: "create player without category defaults to Unranked",
			payloadRequest: request.NewPlayerRequest{
				Sport: "Football",
				ID:    "player4",
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]domain.Entity{}, nil)
				playerRepository.On("Save", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.ID == "player4" && entity.Category == common.Unranked
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
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

