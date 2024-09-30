package team

import (
	"testing"
)

func TestTeamCreationHandler(t *testing.T) {
	// Mocking the validator and use case
	/*validator := validator.New()
	teamRepositoryMock := new(mocks.Repository)
	createTeamUC := usecases.NewCreateTeamUC(teamRepositoryMock)

	controller := team.NewController(createTeamUC, validator)

	// Create a Gin recorder and context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/team", controller.TeamCreationHandler)

	// Mock input
	newTeam := request2.NewTeamRequest{
		Name: "Test Team",
		// other fields...
	}
	jsonData, _ := json.Marshal(newTeam)

	req, _ := http.NewRequest("POST", "/team", bytes.NewBuffer(jsonData))
	resp := httptest.NewRecorder()

	// Define the behavior of the mocked team repository
	teamEntity := team.Entity{Name: "Test Team"} // Adapt fields accordingly
	teamRepositoryMock.On("Save", mock.Anything).Return(nil) // Simulate no error

	// Perform the test
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code, "Expected HTTP status 200 OK")
	assert.NoError(t, json.Unmarshal(resp.Body.Bytes(), &teamEntity), "Expected no error in unmarshalling response")

	// Check that all expected methods were called
	teamRepositoryMock.AssertExpectations(t) */
}
