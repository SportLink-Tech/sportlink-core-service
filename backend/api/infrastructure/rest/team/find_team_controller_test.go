package team_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	team2 "sportlink/api/domain/team"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/team"
	amocks "sportlink/mocks/api/application"
)

func TestFindTeam(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name          string
		sport         string
		nameQuery     string
		categoryQuery string
		given         func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity])
		then          func(t *testing.T, responseCode int, response interface{})
	}{
		{
			name:      "given use case returns teams when find then returns ok",
			sport:     "Football",
			nameQuery: "Boca",
			given: func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {
				slice := []team2.Entity{
					{Name: "Boca Juniors", Sport: common.Football, Category: common.L1, Stats: *common.NewStats(10, 5, 2), Members: []player.Entity{}},
					{Name: "Boca Unidos", Sport: common.Football, Category: common.L2, Stats: *common.NewStats(5, 3, 1), Members: []player.Entity{}},
				}
				find.On("Invoke", mock.Anything, mock.MatchedBy(func(q team2.DomainQuery) bool {
					return len(q.Sports) == 1 && q.Sports[0] == common.Football && q.Name == "Boca"
				})).Return(&slice, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				teams := response.([]interface{})
				assert.Len(t, teams, 2)
			},
		},
		{
			name:      "given no teams when find then returns not found",
			sport:     "Tennis",
			nameQuery: "NonExistent",
			given: func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {
				slice := []team2.Entity{}
				find.On("Invoke", mock.Anything, mock.Anything).Return(&slice, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusNotFound, responseCode)
				m := response.(map[string]interface{})
				assert.Equal(t, "not_found", m["code"])
			},
		},
		{
			name:      "given use case error when find then returns conflict",
			sport:     "Football",
			nameQuery: "River",
			given: func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {
				find.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				m := response.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", m["code"])
			},
		},
		{
			name:      "given missing sport path when find then returns bad request",
			sport:     "",
			nameQuery: "Boca",
			given:     func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				m := response.(map[string]interface{})
				assert.Equal(t, "invalid_request_format", m["code"])
			},
		},
		{
			name:          "given invalid category format when find then returns validation error",
			sport:         "Football",
			categoryQuery: "invalid",
			given:         func(t *testing.T, find *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				m := response.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			createUC := amocks.NewUseCase[team2.Entity, team2.Entity](t)
			retrieveUC := amocks.NewUseCase[team2.ID, team2.Entity](t)
			findUC := amocks.NewUseCase[team2.DomainQuery, []team2.Entity](t)
			listUC := amocks.NewUseCase[team2.DomainQuery, []team2.Entity](t)
			tc.given(t, findUC)

			controller := team.NewController(createUC, retrieveUC, findUC, listUC, v)
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.GET("/sport/:sport/team", controller.FindTeam)

			url := fmt.Sprintf("/sport/%s/team", tc.sport)
			q := []string{}
			if tc.nameQuery != "" {
				q = append(q, "name="+tc.nameQuery)
			}
			if tc.categoryQuery != "" {
				q = append(q, "category="+tc.categoryQuery)
			}
			if len(q) > 0 {
				url += "?" + q[0]
				for i := 1; i < len(q); i++ {
					url += "&" + q[i]
				}
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var response interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &response)
			tc.then(t, rec.Code, response)
		})
	}
}
