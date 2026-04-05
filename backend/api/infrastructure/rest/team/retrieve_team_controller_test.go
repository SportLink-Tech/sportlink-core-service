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

func TestRetrieveTeam(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name   string
		sport  string
		tname  string
		given  func(t *testing.T, retrieve *amocks.UseCase[team2.ID, team2.Entity])
		then   func(t *testing.T, code int, body interface{})
	}{
		{
			name:  "given use case returns team when retrieve then returns ok",
			sport: "Football",
			tname: "Boca",
			given: func(t *testing.T, retrieve *amocks.UseCase[team2.ID, team2.Entity]) {
				ent := team2.NewTeam("Boca", common.L1, *common.NewStats(1, 0, 0), common.Football, []player.Entity{}, "acc")
				retrieve.On("Invoke", mock.Anything, mock.MatchedBy(func(id team2.ID) bool {
					return id.Name == "Boca" && id.Sport == common.Football
				})).Return(&ent, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "Boca", m["Name"])
			},
		},
		{
			name:  "given use case error when retrieve then returns conflict",
			sport: "Paddle",
			tname: "Missing",
			given: func(t *testing.T, retrieve *amocks.UseCase[team2.ID, team2.Entity]) {
				retrieve.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
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
			tc.given(t, retrieveUC)

			ctl := team.NewController(createUC, retrieveUC, findUC, listUC, v)
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.GET("/sport/:sport/team/:team", ctl.RetrieveTeam)

			req := httptest.NewRequest(http.MethodGet, "/sport/"+tc.sport+"/team/"+tc.tname, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &body)
			tc.then(t, rec.Code, body)
		})
	}
}
