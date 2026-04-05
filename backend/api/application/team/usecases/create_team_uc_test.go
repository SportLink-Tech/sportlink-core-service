package usecases_test

import (
	"context"
	"errors"
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
	"github.com/stretchr/testify/require"
)

func TestNewCreateTeamUC(t *testing.T) {
	ctx := context.Background()
	findErr := errors.New("player lookup failed")
	saveErr := errors.New("it was an error")

	tests := []struct {
		name   string
		entity team.Entity
		given  func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository)
		then   func(t *testing.T, response *team.Entity, err error)
	}{
		{
			name: "given team has no members and save succeeds when invoke then returns team",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				make([]player.Entity, 0),
				"",
			),
			given: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(tm team.Entity) bool {
					return tm.Name == "Boca Jr" &&
						tm.Category == common.L1 &&
						tm.Sport == common.Football &&
						len(tm.Members) == 0 &&
						tm.Stats == *common.NewStats(10, 0, 0)
				})).Return(nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				require.NoError(t, err)
				expected := team.NewTeam(
					"Boca Jr",
					common.L1,
					*common.NewStats(10, 0, 0),
					common.Football,
					make([]player.Entity, 0),
					"",
				)
				require.NotNil(t, response)
				assert.Equal(t, expected.ID, response.ID)
				assert.Equal(t, expected.Name, response.Name)
				assert.Equal(t, expected.Category, response.Category)
				assert.Equal(t, expected.Sport, response.Sport)
			},
		},
		{
			name: "given all members exist when invoke then saves and returns team",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				[]player.Entity{
					{ID: "eldiegote", Category: common.L1, Sport: common.Football},
					{ID: "elpajaro", Category: common.L1, Sport: common.Football},
				},
				"",
			),
			given: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.MatchedBy(func(query player.DomainQuery) bool {
					return reflect.DeepEqual(query.Ids, []string{"eldiegote", "elpajaro"})
				})).Return([]player.Entity{
					{ID: "eldiegote", Category: common.L1, Sport: common.Football},
					{ID: "elpajaro", Category: common.L1, Sport: common.Football},
				}, nil)

				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(tm team.Entity) bool {
					return tm.Name == "Boca Jr" &&
						len(tm.Members) == 2 &&
						tm.Stats == *common.NewStats(10, 0, 0)
				})).Return(nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
				assert.Len(t, response.Members, 2)
			},
		},
		{
			name: "given not all members exist when invoke then returns error and does not save",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				[]player.Entity{
					{ID: "eldiegote", Category: common.L1, Sport: common.Football},
					{ID: "elpajaro", Category: common.L1, Sport: common.Football},
				},
				"",
			),
			given: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.MatchedBy(func(query player.DomainQuery) bool {
					return reflect.DeepEqual(query.Ids, []string{"eldiegote", "elpajaro"})
				})).Return([]player.Entity{
					{ID: "elpajaro", Category: common.L1, Sport: common.Football},
				}, nil)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, response)
				assert.Contains(t, err.Error(), "some of the team member does not exist")
			},
		},
		{
			name: "given player find fails when invoke then returns wrapped error and does not save",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				[]player.Entity{{ID: "only-one", Category: common.L1, Sport: common.Football}},
				"",
			),
			given: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				playerRepository.On("Find", mock.Anything, mock.Anything).Return(nil, findErr)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, response)
				assert.Contains(t, err.Error(), "error while finding players")
				assert.ErrorIs(t, err, findErr)
			},
		},
		{
			name: "given validation passes but save fails when invoke then returns wrapped error",
			entity: team.NewTeam(
				"Boca Jr",
				common.L1,
				*common.NewStats(10, 0, 0),
				common.Football,
				make([]player.Entity, 0),
				"",
			),
			given: func(t *testing.T, playerRepository *pmocks.Repository, teamRepository *mmocks.Repository) {
				teamRepository.On("Save", mock.Anything, mock.MatchedBy(func(tm team.Entity) bool {
					return tm.Name == "Boca Jr" && tm.Sport == common.Football && len(tm.Members) == 0
				})).Return(saveErr)
			},
			then: func(t *testing.T, response *team.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, response)
				assert.Contains(t, err.Error(), "error while inserting team in database")
				assert.ErrorIs(t, err, saveErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			playerRepository := &pmocks.Repository{}
			teamRepository := &mmocks.Repository{}
			uc := usecases.NewCreateTeamUC(playerRepository, teamRepository)

			tt.given(t, playerRepository, teamRepository)

			response, err := uc.Invoke(ctx, tt.entity)

			tt.then(t, response, err)
			playerRepository.AssertExpectations(t)
			teamRepository.AssertExpectations(t)
		})
	}
}
