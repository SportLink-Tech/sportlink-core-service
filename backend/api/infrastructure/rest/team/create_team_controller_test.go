package team_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	request2 "sportlink/api/application/team/request"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	team2 "sportlink/api/domain/team"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/team"
	amocks "sportlink/mocks/api/application"
)

func TestCreateTeam(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name           string
		payloadRequest request2.NewTeamRequest
		given          func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity])
		then           func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "given use case succeeds when creating team with players then returns created",
			payloadRequest: request2.NewTeamRequest{
				Sport:     "Football",
				Name:      "Boca Juniors",
				Category:  1,
				PlayerIds: []string{"1", "2"},
			},
			given: func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {
				out := team2.NewTeam("Boca Juniors", common.L1, *common.NewStats(0, 0, 0), common.Football,
					[]player.Entity{{ID: "1"}, {ID: "2"}}, "test-account")
				create.On("Invoke", mock.Anything, mock.MatchedBy(func(e team2.Entity) bool {
					return e.Name == "Boca Juniors" && e.Sport == common.Football && e.OwnerAccountID == "test-account"
				})).Return(&out, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "Boca Juniors", response["Name"])
			},
		},
		{
			name: "given use case fails when some players missing then returns conflict",
			payloadRequest: request2.NewTeamRequest{
				Sport:     "Football",
				Name:      "Boca Juniors",
				Category:  1,
				PlayerIds: []string{"1", "2"},
			},
			given: func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {
				create.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
		{
			name: "given use case succeeds when creating team without players then returns created",
			payloadRequest: request2.NewTeamRequest{
				Sport:    "Football",
				Name:     "Boca Juniors",
				Category: 1,
			},
			given: func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {
				out := team2.NewTeam("Boca Juniors", common.L1, *common.NewStats(0, 0, 0), common.Football, nil, "test-account")
				create.On("Invoke", mock.Anything, mock.MatchedBy(func(e team2.Entity) bool {
					return e.Name == "Boca Juniors" && len(e.Members) == 0
				})).Return(&out, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
			},
		},
		{
			name:           "given invalid category when creating then returns bad request",
			payloadRequest: request2.NewTeamRequest{Sport: "Football", Name: "Boca Juniors", Category: 9},
			given:          func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Contains(t, response["message"], "invalid category value")
			},
		},
		{
			name:           "given invalid sport when creating then returns bad request",
			payloadRequest: request2.NewTeamRequest{Sport: "fuchibol", Name: "River Plate", Category: 2},
			given:          func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Contains(t, response["message"], "Sport")
			},
		},
		{
			name:           "given empty name when creating then returns bad request",
			payloadRequest: request2.NewTeamRequest{Sport: "Football", Name: "", Category: 1},
			given:          func(t *testing.T, create *amocks.UseCase[team2.Entity, team2.Entity]) {},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Contains(t, response["message"], "Name")
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
			tc.given(t, createUC)

			controller := team.NewController(createUC, retrieveUC, findUC, listUC, v)
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.POST("/account/:accountId/team", controller.CreateTeam)

			jsonData, _ := json.Marshal(tc.payloadRequest)
			req, _ := http.NewRequest(http.MethodPost, "/account/test-account/team", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var response map[string]interface{}
			_ = json.Unmarshal(rec.Body.Bytes(), &response)
			tc.then(t, rec.Code, response)
		})
	}
}
