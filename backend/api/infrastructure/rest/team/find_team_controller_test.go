package team_test

import (
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

func TestFindTeamController(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name          string
		sport         string
		nameQuery     string
		categoryQuery string
		on            func(t *testing.T, teamRepository *tmocks.Repository)
		assertions    func(t *testing.T, responseCode int, response interface{})
	}{
		{
			name:      "find teams successfully - multiple results",
			sport:     "Football",
			nameQuery: "Boca",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "Boca"
				})).Return([]team2.Entity{
					{
						Name:     "Boca Juniors",
						Sport:    common.Football,
						Category: common.L1,
						Stats:    *common.NewStats(10, 5, 2),
						Members:  []player.Entity{},
					},
					{
						Name:     "Boca Unidos",
						Sport:    common.Football,
						Category: common.L2,
						Stats:    *common.NewStats(5, 3, 1),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 2)
				firstTeam := teams[0].(map[string]interface{})
				assert.Equal(t, "Boca Juniors", firstTeam["Name"])
				secondTeam := teams[1].(map[string]interface{})
				assert.Equal(t, "Boca Unidos", secondTeam["Name"])
			},
		},
		{
			name:      "find teams successfully - single result",
			sport:     "Paddle",
			nameQuery: "Los Delfines",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						query.Name == "Los Delfines"
				})).Return([]team2.Entity{
					{
						Name:     "Los Delfines",
						Sport:    common.Paddle,
						Category: common.L7,
						Stats:    *common.NewStats(15, 2, 1),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 1)
				firstTeam := teams[0].(map[string]interface{})
				assert.Equal(t, "Los Delfines", firstTeam["Name"])
			},
		},
		{
			name:      "find teams returns 404 - no results",
			sport:     "Tennis",
			nameQuery: "NonExistent",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Tennis &&
						query.Name == "NonExistent"
				})).Return([]team2.Entity{}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusNotFound, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "not_found", responseMap["code"])
				assert.Contains(t, responseMap["message"], "No teams found")
			},
		},
		{
			name:      "find teams fails - repository error",
			sport:     "Football",
			nameQuery: "River",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "River"
				})).Return([]team2.Entity{}, fmt.Errorf("database connection error"))
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", responseMap["code"])
				assert.Contains(t, responseMap["message"], "database connection error")
			},
		},
		{
			name:      "find teams fails - missing sport parameter",
			sport:     "",
			nameQuery: "Boca",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				// No mock setup needed as validation fails before repository call
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "invalid_request_format", responseMap["code"])
			},
		},
		{
			name:          "find teams by category only",
			sport:         "Football",
			nameQuery:     "",
			categoryQuery: "1",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						len(query.Categories) == 1 &&
						query.Categories[0] == common.L1 &&
						query.Name == ""
				})).Return([]team2.Entity{
					{
						Name:     "Boca Juniors",
						Sport:    common.Football,
						Category: common.L1,
						Stats:    *common.NewStats(10, 5, 2),
						Members:  []player.Entity{},
					},
					{
						Name:     "River Plate",
						Sport:    common.Football,
						Category: common.L1,
						Stats:    *common.NewStats(8, 3, 1),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 2)
			},
		},
		{
			name:          "find teams by multiple categories",
			sport:         "Paddle",
			nameQuery:     "",
			categoryQuery: "5,7",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						len(query.Categories) == 2 &&
						query.Categories[0] == common.L5 &&
						query.Categories[1] == common.L7
				})).Return([]team2.Entity{
					{
						Name:     "Los Tiburones",
						Sport:    common.Paddle,
						Category: common.L5,
						Stats:    *common.NewStats(12, 2, 0),
						Members:  []player.Entity{},
					},
					{
						Name:     "Los Delfines",
						Sport:    common.Paddle,
						Category: common.L7,
						Stats:    *common.NewStats(15, 1, 0),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 2)
			},
		},
		{
			name:          "find teams by name and category",
			sport:         "Football",
			nameQuery:     "Boca",
			categoryQuery: "1",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "Boca" &&
						len(query.Categories) == 1 &&
						query.Categories[0] == common.L1
				})).Return([]team2.Entity{
					{
						Name:     "Boca Juniors",
						Sport:    common.Football,
						Category: common.L1,
						Stats:    *common.NewStats(10, 5, 2),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 1)
				firstTeam := teams[0].(map[string]interface{})
				assert.Equal(t, "Boca Juniors", firstTeam["Name"])
			},
		},
		{
			name:          "find all teams by sport only (no name, no category)",
			sport:         "Paddle",
			nameQuery:     "",
			categoryQuery: "",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(query team2.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						query.Name == "" &&
						len(query.Categories) == 0
				})).Return([]team2.Entity{
					{
						Name:     "Los Delfines",
						Sport:    common.Paddle,
						Category: common.L7,
						Stats:    *common.NewStats(15, 2, 1),
						Members:  []player.Entity{},
					},
					{
						Name:     "Los Tiburones",
						Sport:    common.Paddle,
						Category: common.L5,
						Stats:    *common.NewStats(12, 3, 0),
						Members:  []player.Entity{},
					},
					{
						Name:     "Los Orcas",
						Sport:    common.Paddle,
						Category: common.L3,
						Stats:    *common.NewStats(8, 4, 2),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 3)
				// Verify we got teams from different categories
				teamNames := make([]string, len(teams))
				for i, team := range teams {
					teamMap := team.(map[string]interface{})
					teamNames[i] = teamMap["Name"].(string)
					assert.Equal(t, "Paddle", teamMap["Sport"])
				}
				assert.Contains(t, teamNames, "Los Delfines")
				assert.Contains(t, teamNames, "Los Tiburones")
				assert.Contains(t, teamNames, "Los Orcas")
			},
		},
		{
			name:          "find teams fails - invalid category format",
			sport:         "Football",
			nameQuery:     "",
			categoryQuery: "invalid",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				// No mock setup needed as validation fails before repository call
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", responseMap["code"])
				assert.Contains(t, responseMap["message"], "invalid category format")
			},
		},
		{
			name:          "find teams fails - invalid category value",
			sport:         "Football",
			nameQuery:     "",
			categoryQuery: "99",
			on: func(t *testing.T, teamRepository *tmocks.Repository) {
				// No mock setup needed as validation fails before repository call
			},
			assertions: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", responseMap["code"])
				assert.Contains(t, responseMap["message"], "invalid category value")
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

			controller := team.NewController(createTeamUC, retrieveTeamUC, findTeamUC, validator)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.GET("/sport/:sport/team", controller.FindTeam)

			// given
			tc.on(t, teamRepository)
			url := fmt.Sprintf("/sport/%s/team", tc.sport)

			// Build query string
			queryParams := []string{}
			if tc.nameQuery != "" {
				queryParams = append(queryParams, fmt.Sprintf("name=%s", tc.nameQuery))
			}
			if tc.categoryQuery != "" {
				queryParams = append(queryParams, fmt.Sprintf("category=%s", tc.categoryQuery))
			}
			if len(queryParams) > 0 {
				url += "?" + queryParams[0]
				for i := 1; i < len(queryParams); i++ {
					url += "&" + queryParams[i]
				}
			}

			req, _ := http.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()

			// when
			router.ServeHTTP(resp, req)

			// then
			var response interface{}
			json.Unmarshal(resp.Body.Bytes(), &response)
			tc.assertions(t, resp.Code, response)
		})
	}
}
