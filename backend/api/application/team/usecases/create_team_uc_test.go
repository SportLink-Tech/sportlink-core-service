package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
	pmocks "sportlink/mocks/api/domain/player"
	mmocks "sportlink/mocks/api/domain/team"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTeamUC_Invoke(t *testing.T) {

	tests := []struct {
		name   string
		entity team.Entity
		on     func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository)
		then   func(t *testing.T, response *team.Entity, err error)
	}{
		{
			name: "save team successfully",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				make([]player.Entity, 0),
			),
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(team team.Entity) bool {
					return team.Name == "Boca Jr" &&
						team.Category == common.L1 &&
						team.Sport == common.Football &&
						len(team.Members) == 0 &&
						team.Stats == *common.NewStats(10, 0, 0)
				})).Return(nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				assert.NoError(t, err)
				expected := team.NewTeam(
					"Boca Jr",
					common.L1,
					*common.NewStats(10, 0, 0),
					common.Football,
					make([]player.Entity, 0),
				)
				assert.Equal(t, expected.ID, response.ID)
				assert.Equal(t, expected.Name, response.Name)
				assert.Equal(t, expected.Category, response.Category)
				assert.Equal(t, expected.Sport, response.Sport)
			},
		},
		{
			name: "save team with players successfully",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				[]player.Entity{
					{
						ID:       "eldiegote",
						Category: common.L1,
						Sport:    common.Football,
					},
					{
						ID:       "elpajaro",
						Category: common.L1,
						Sport:    common.Football,
					},
				},
			),
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(team team.Entity) bool {
					return team.Name == "Boca Jr" &&
						team.Category == common.L1 &&
						team.Sport == common.Football &&
						len(team.Members) == 2 &&
						team.Stats == *common.NewStats(10, 0, 0)
				})).Return(nil)

				playerRepository.On("Find", mock.Anything, mock.MatchedBy(func(query player.DomainQuery) bool {
					return reflect.DeepEqual(query.Ids, []string{"eldiegote", "elpajaro"})
				})).Return([]player.Entity{
					{
						ID:       "eldiegote",
						Category: common.L1,
						Sport:    common.Football,
					},
					{
						ID:       "elpajaro",
						Category: common.L1,
						Sport:    common.Football,
					},
				}, nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				assert.NoError(t, err)
				expected := team.NewTeam(
					"Boca Jr",
					common.L1,
					*common.NewStats(10, 0, 0),
					common.Football,
					[]player.Entity{
						{
							ID:       "eldiegote",
							Category: common.L1,
							Sport:    common.Football,
						},
						{
							ID:       "elpajaro",
							Category: common.L1,
							Sport:    common.Football,
						},
					},
				)
				assert.Equal(t, expected.ID, response.ID)
				assert.Equal(t, expected.Name, response.Name)
				assert.Equal(t, expected.Category, response.Category)
				assert.Equal(t, expected.Sport, response.Sport)
				assert.Equal(t, len(expected.Members), len(response.Members))
			},
		},
		{
			name: "when some of the team players does not exist then the team could not be created",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				[]player.Entity{
					{
						ID:       "eldiegote",
						Category: common.L1,
						Sport:    common.Football,
					},
					{
						ID:       "elpajaro",
						Category: common.L1,
						Sport:    common.Football,
					},
				},
			),
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(team team.Entity) bool {
					return team.Name == "Boca Jr" &&
						team.Category == common.L1 &&
						team.Sport == common.Football &&
						len(team.Members) == 2 &&
						team.Stats == *common.NewStats(10, 0, 0)
				})).Return(nil)

				playerRepository.On("Find", mock.Anything, mock.MatchedBy(func(query player.DomainQuery) bool {
					return reflect.DeepEqual(query.Ids, []string{"eldiegote", "elpajaro"})
				})).Return([]player.Entity{
					{
						ID:       "elpajaro",
						Category: common.L1,
						Sport:    common.Football,
					},
				}, nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "some of the team member does not exist")
				assert.Nil(t, response)
			},
		},
		{
			name: "when the save method repository thrown an error, then it must be retrieved",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				make([]player.Entity, 0),
			),
			on: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(team team.Entity) bool {
					return team.Name == "Boca Jr" &&
						team.Category == common.L1 &&
						team.Sport == common.Football &&
						len(team.Members) == 0 &&
						team.Stats == *common.NewStats(10, 0, 0)
				})).Return(fmt.Errorf("it was an error"))
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error while inserting team in database")
				assert.Nil(t, response)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			playerRepository := &pmocks.Repository{}
			teamRepository := &mmocks.Repository{}
			uc := usecases.NewCreateTeamUC(playerRepository, teamRepository)

			// given
			tt.on(t, playerRepository, teamRepository)

			// when
			response, err := uc.Invoke(context.Background(), tt.entity)

			// then
			tt.then(t, response, err)
		})
	}
}
