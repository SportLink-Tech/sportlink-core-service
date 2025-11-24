package usecases_test

import (
	"fmt"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
	mmocks "sportlink/mocks/api/domain/team"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindTeamUC_Invoke(t *testing.T) {

	tests := []struct {
		name  string
		query team.DomainQuery
		on    func(t *testing.T, repository *mmocks.Repository)
		then  func(t *testing.T, result *[]team.Entity, err error)
	}{
		{
			name: "find teams successfully - multiple results",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Football},
				Name:   "Boca",
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "Boca"
				})).Return([]team.Entity{
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
			then: func(t *testing.T, result *[]team.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 2)
				assert.Equal(t, "Boca Juniors", (*result)[0].Name)
				assert.Equal(t, "Boca Unidos", (*result)[1].Name)
			},
		},
		{
			name: "find teams successfully - single result",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Name:   "Los Delfines",
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						query.Name == "Los Delfines"
				})).Return([]team.Entity{
					{
						Name:     "Los Delfines",
						Sport:    common.Paddle,
						Category: common.L7,
						Stats:    *common.NewStats(15, 2, 1),
						Members:  []player.Entity{},
					},
				}, nil)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "Los Delfines", (*result)[0].Name)
				assert.Equal(t, common.Paddle, (*result)[0].Sport)
			},
		},
		{
			name: "find teams successfully - no results",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Tennis},
				Name:   "NonExistent",
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Tennis &&
						query.Name == "NonExistent"
				})).Return([]team.Entity{}, nil)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Empty(t, *result)
			},
		},
		{
			name: "find teams fails - repository error",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Football},
				Name:   "River",
			},
			on: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "River"
				})).Return([]team.Entity{}, fmt.Errorf("database connection error"))
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database connection error", err.Error())
				assert.Nil(t, result)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			teamRepository := &mmocks.Repository{}
			uc := usecases.NewFindTeamUC(teamRepository)

			// given
			tt.on(t, teamRepository)

			// when
			result, err := uc.Invoke(tt.query)

			// then
			tt.then(t, result, err)
		})
	}

}
