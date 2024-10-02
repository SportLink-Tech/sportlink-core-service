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
	router.POST("/team", controller.TeamCreationHandler)

	// Crear varios casos de prueba para diferentes combinaciones de campos vacíos
	testCases := []struct {
		description  string
		payload      request2.NewTeamRequest
		expectedCode int
	}{
		{
			description:  "empty name",
			payload:      request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 1},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			teamRepository.On("Save", mock.Anything).Return(nil)

			jsonData, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/team", bytes.NewBuffer(jsonData))
			resp := httptest.NewRecorder()

			// Perform the test
			router.ServeHTTP(resp, req)

			// Assertions
			assert.Equal(t, tc.expectedCode, resp.Code, "Expected HTTP status matches the expected")
		})
	}
}
