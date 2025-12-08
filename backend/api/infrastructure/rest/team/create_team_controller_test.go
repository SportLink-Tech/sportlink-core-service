package team_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	request2 "sportlink/api/application/team/request"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	team2 "sportlink/api/domain/team"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/team"
	pmocks "sportlink/mocks/api/domain/player"
	tmocks "sportlink/mocks/api/domain/team"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTeamCreationHandlerWithEmptyFields(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name           string
		payloadRequest request2.NewTeamRequest
		on             func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository)
		assertions     func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "create a new Boca Juniors team with player ids successfully",
			payloadRequest: request2.NewTeamRequest{
				Sport:     "Football",
				Name:      "Boca Juniors",
				Category:  1,
				PlayerIds: []string{"1", "2"},
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(entity team2.Entity) bool {
					return entity.Sport == common.Football &&
						entity.Stats == *common.NewStats(0, 0, 0) &&
						entity.Name == "Boca Juniors" &&
						entity.Category == common.L1
				})).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{
					{
						ID: "1",
					},
					{
						ID: "2",
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name: "team creation failed when some of the players does not exist",
			payloadRequest: request2.NewTeamRequest{
				Sport:     "Football",
				Name:      "Boca Juniors",
				Category:  1,
				PlayerIds: []string{"1", "2"},
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{
					{
						ID: "1",
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, response["code"], "use_case_execution_error")
				assert.Equal(t, response["message"], "use case execution failed. Err: some of the team member does not exist")
			},
		},
		{
			name: "create a new Boca Juniors team without players successfully",
			payloadRequest: request2.NewTeamRequest{
				Sport:    "Football",
				Name:     "Boca Juniors",
				Category: 1,
			},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "create a new River Plate team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "Football", Name: "River Plate", Category: 2},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "create a new Los Delfines paddle team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "Paddle", Name: "Los Delfines", Category: 7},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "fails when create a new team with invalid category",
			payloadRequest: request2.NewTeamRequest{Sport: "Football", Name: "Boca Juniors", Category: 9},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Err: invalid category value: 9")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name:           "fails when create a new team with invalid sport",
			payloadRequest: request2.NewTeamRequest{Sport: "fuchibol", Name: "River Plate", Category: 2},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Error:Field validation for 'Sport' failed")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name:           "create a new team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "Football", Name: "", Category: 1},
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *tmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
				playerRepository.On("Find", mock.Anything, mock.Anything).Return([]player.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], " Error:Field validation for 'Name'")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			playerRepository := new(pmocks.Repository)
			teamRepository := new(tmocks.Repository)
			createTeamUC := usecases.NewCreateTeamUC(playerRepository, teamRepository)
			retrieveTeamUC := usecases.NewRetrieveTeamUC(teamRepository)
			findTeamUC := usecases.NewFindTeamUC(teamRepository)

			controller := team.NewController(createTeamUC, retrieveTeamUC, findTeamUC, validator)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.POST("/team", controller.CreateTeam)

			// given
			tc.on(t, playerRepository, teamRepository)
			jsonData, _ := json.Marshal(tc.payloadRequest)
			req, _ := http.NewRequest("POST", "/team", bytes.NewBuffer(jsonData))
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
