package team_test

import (
	"encoding/json"
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

func TestListAccountTeams(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name  string
		accID string
		given func(t *testing.T, list *amocks.UseCase[team2.DomainQuery, []team2.Entity])
		then  func(t *testing.T, code int, body interface{})
	}{
		{
			name:  "given use case returns teams when list by account then returns ok",
			accID: "owner-1",
			given: func(t *testing.T, list *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {
				slice := []team2.Entity{
					{Name: "A", Sport: common.Paddle, Category: common.L5, Stats: *common.NewStats(0, 0, 0), Members: []player.Entity{}, OwnerAccountID: "owner-1"},
				}
				list.On("Invoke", mock.Anything, mock.MatchedBy(func(q team2.DomainQuery) bool {
					return q.OwnerAccountID == "owner-1"
				})).Return(&slice, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
			},
		},
		{
			name:  "given use case fails when list then returns conflict",
			accID: "owner-1",
			given: func(t *testing.T, list *amocks.UseCase[team2.DomainQuery, []team2.Entity]) {
				list.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusConflict, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", m["code"])
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
			tc.given(t, listUC)

			ctl := team.NewController(createUC, retrieveUC, findUC, listUC, v)
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.GET("/account/:accountId/team", ctl.ListAccountTeams)

			req := httptest.NewRequest(http.MethodGet, "/account/"+tc.accID+"/team", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &body)
			tc.then(t, rec.Code, body)
		})
	}
}
