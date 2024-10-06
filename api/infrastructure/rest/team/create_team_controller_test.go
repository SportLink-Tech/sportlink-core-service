package team

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
	"testing"
)

func TestTeamCreationHandlerWithEmptyFields(t *testing.T) {
	validator := validator.New()

	playerRepository := new(pmocks.Repository)
	teamRepository := new(tmocks.Repository)
	createTeamUC := usecases.NewCreateTeamUC(playerRepository, teamRepository)

	controller := NewController(createTeamUC, validator)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/team", controller.TeamCreationHandler)

	testCases := []struct {
		name           string
		payloadRequest request2.NewTeamRequest
		expectedCode   int
	}{
		{
			name:           "create a new team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 1},
			expectedCode:   http.StatusCreated,
		},
		{
			name:           "fails when create a new team with invalid category",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 9},
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "fails when create a new team with invalid sport",
			payloadRequest: request2.NewTeamRequest{Sport: "fuchibol", Name: "River Plate", Category: 2},
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "create a new team successfully",
			payloadRequest: request2.NewTeamRequest{Sport: "football", Name: "", Category: 1},
			expectedCode:   http.StatusBadRequest,
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
			var response map[string]string
			json.Unmarshal(resp.Body.Bytes(), &response)
			assert.Contains(t, response["message"], "Err: Key: 'NewTeamRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag")
			assert.Equal(t, tc.expectedCode, resp.Code, "Expected HTTP status matches the expected")
		})
	}
}
