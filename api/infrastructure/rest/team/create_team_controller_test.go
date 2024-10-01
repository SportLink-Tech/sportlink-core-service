package team

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	request2 "sportlink/api/application/team/request"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/player/mocks"
	"testing"
)

func TestTeamCreationHandlerWithEmptyFields(t *testing.T) {
	validator := validator.New()

	playerRepositoryMock := new(mocks.Repository)
	createTeamUC := usecases.NewCreateTeamUC(playerRepositoryMock)

	controller := NewController(createTeamUC, validator)

	// Create a Gin recorder and context
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
		{
			description:  "empty sport",
			payload:      request2.NewTeamRequest{Sport: "", Name: "Boca Juniors", Category: 1},
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "empty category",
			payload:      request2.NewTeamRequest{Sport: "football", Name: "Boca Juniors", Category: 0}, // Assuming Category cannot be zero if required
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
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
