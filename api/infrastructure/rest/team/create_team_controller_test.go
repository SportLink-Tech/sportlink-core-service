package team_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	request2 "sportlink/api/application/team/request"
	"sportlink/api/application/team/usecases"
	pmocks "sportlink/api/domain/player/mocks"
	tmocks "sportlink/api/domain/team/mocks"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/team"
	"testing"
)

func TestTeamCreationHandlerWithEmptyFields(t *testing.T) {
	validator := validator.New()

	playerRepository := new(pmocks.Repository)
	teamRepository := new(tmocks.Repository)
	createTeamUC := usecases.NewCreateTeamUC(playerRepository, teamRepository)

	controller := team.NewController(createTeamUC, validator)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/team", controller.TeamCreationHandler)

	testCases := []struct {
		name           string
		payloadRequest request2.NewTeamRequest
		assertions     func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "create a new Boca Juniors team with player ids successfully",
			payloadRequest: request2.NewTeamRequest{
				Sport:     "football",
				Name:      "Boca Juniors",
				Category:  1,
				PlayerIds: []string{"1", "2"},
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "create a new Boca Juniors team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 1},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "create a new River Plate team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "River Plate", Category: 2},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "create a new Los Delfines paddle team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "paddle", Name: "Los Delfines", Category: 7},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "fails when create a new team with invalid category",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 9},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Err: invalid category value: 9")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name:           "fails when create a new team with invalid sport",
			payloadRequest: request2.NewTeamRequest{Sport: "fuchibol", Name: "River Plate", Category: 2},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], "Error:Field validation for 'Sport' failed")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
		{
			name:           "create a new team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "", Category: 1},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Contains(t, response["message"], " Error:Field validation for 'Name'")
				assert.Equal(t, http.StatusBadRequest, responseCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			teamRepository.On("Save", mock.Anything).Return(nil)

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
