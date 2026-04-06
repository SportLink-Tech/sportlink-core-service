package team_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	team2 "sportlink/api/domain/team"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/team"
	tmocks "sportlink/mocks/api/domain/team"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTeamController(t *testing.T) {
	v := validator.New()

	existingTeam := team2.Entity{
		ID:       "SPORT#Football#NAME#Boca Juniors",
		Name:     "Boca Juniors",
		Sport:    common.Football,
		Category: common.L1,
		Stats:    *common.NewStats(0, 0, 0),
		Members:  []player.Entity{},
	}

	testCases := []struct {
		name       string
		sport      string
		teamName   string
		body       map[string]interface{}
		on         func(t *testing.T, teamRepository *tmocks.Repository)
		assertions func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name:     "updates team name successfully",
			sport:    "Football",
			teamName: "Boca Juniors",
			body:     map[string]interface{}{"name": "Boca Senior"},
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(q team2.DomainQuery) bool {
					return q.Name == "Boca Juniors" && len(q.Sports) == 1 && q.Sports[0] == common.Football
				})).Return([]team2.Entity{existingTeam}, nil)
				teamRepository.On("Update", mock.Anything, "SPORT#Football#NAME#Boca Juniors", mock.MatchedBy(func(e team2.Entity) bool {
					return e.Name == "Boca Senior"
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				assert.Equal(t, "Boca Senior", response["Name"])
			},
		},
		{
			name:     "patch with empty body keeps existing name",
			sport:    "Football",
			teamName: "Boca Juniors",
			body:     map[string]interface{}{},
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.Anything).Return([]team2.Entity{existingTeam}, nil)
				teamRepository.On("Update", mock.Anything, "SPORT#Football#NAME#Boca Juniors", mock.MatchedBy(func(e team2.Entity) bool {
					return e.Name == "Boca Juniors"
				})).Return(nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				assert.Equal(t, "Boca Juniors", response["Name"])
			},
		},
		{
			name:     "fails when team does not exist",
			sport:    "Football",
			teamName: "Unknown Team",
			body:     map[string]interface{}{"name": "New Name"},
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.Anything).Return([]team2.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
				assert.Contains(t, response["message"], "team not found")
			},
		},
		{
			name:     "fails when repository returns error",
			sport:    "Football",
			teamName: "Boca Juniors",
			body:     map[string]interface{}{"name": "Boca Senior"},
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.Anything).Return([]team2.Entity{}, fmt.Errorf("db error"))
			},
			assertions: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			teamRepository := new(tmocks.Repository)
			createTeamUC := usecases.NewCreateTeamUC(nil, teamRepository)
			retrieveTeamUC := usecases.NewRetrieveTeamUC(teamRepository)
			findTeamUC := usecases.NewFindTeamUC(teamRepository)
			updateTeamUC := usecases.NewUpdateTeamUC(teamRepository)

			controller := team.NewController(createTeamUC, retrieveTeamUC, findTeamUC, findTeamUC, updateTeamUC, v)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.PATCH("/sport/:sport/team/:team", controller.UpdateTeam)

			tc.on(t, teamRepository)
			jsonData, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/sport/%s/team/%s", tc.sport, tc.teamName), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			response := createMapResponse(resp)
			tc.assertions(t, resp.Code, response)
		})
	}
}
