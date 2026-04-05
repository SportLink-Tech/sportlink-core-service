package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/team/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
	mmocks "sportlink/mocks/api/domain/team"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFindTeamUC(t *testing.T) {
	ctx := context.Background()
	findErr := errors.New("database connection error")

	tests := []struct {
		name  string
		query team.DomainQuery
		given func(t *testing.T, repository *mmocks.Repository)
		then  func(t *testing.T, result *[]team.Entity, err error)
	}{
		{
			name: "given repository returns multiple teams when invoke then returns pointer to slice",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Football},
				Name:   "Boca",
			},
			given: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "Boca"
				})).Return([]team.Entity{
					{
						Name: "Boca Juniors", Sport: common.Football, Category: common.L1,
						Stats: *common.NewStats(10, 5, 2), Members: []player.Entity{},
					},
					{
						Name: "Boca Unidos", Sport: common.Football, Category: common.L2,
						Stats: *common.NewStats(5, 3, 1), Members: []player.Entity{},
					},
				}, nil)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, *result, 2)
				assert.Equal(t, "Boca Juniors", (*result)[0].Name)
				assert.Equal(t, "Boca Unidos", (*result)[1].Name)
			},
		},
		{
			name: "given repository returns single team when invoke then returns one element",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Name:   "Los Delfines",
			},
			given: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						query.Name == "Los Delfines"
				})).Return([]team.Entity{
					{
						Name: "Los Delfines", Sport: common.Paddle, Category: common.L7,
						Stats: *common.NewStats(15, 2, 1), Members: []player.Entity{},
					},
				}, nil)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "Los Delfines", (*result)[0].Name)
				assert.Equal(t, common.Paddle, (*result)[0].Sport)
			},
		},
		{
			name: "given repository returns empty when invoke then returns empty slice pointer",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Tennis},
				Name:   "NonExistent",
			},
			given: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Tennis &&
						query.Name == "NonExistent"
				})).Return([]team.Entity{}, nil)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, *result)
			},
		},
		{
			name: "given repository fails when invoke then returns error",
			query: team.DomainQuery{
				Sports: []common.Sport{common.Football},
				Name:   "River",
			},
			given: func(t *testing.T, repository *mmocks.Repository) {
				repository.On("Find", mock.Anything, mock.MatchedBy(func(query team.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Football &&
						query.Name == "River"
				})).Return([]team.Entity{}, findErr)
			},
			then: func(t *testing.T, result *[]team.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.ErrorIs(t, err, findErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamRepository := &mmocks.Repository{}
			uc := usecases.NewFindTeamUC(teamRepository)

			tt.given(t, teamRepository)

			result, err := uc.Invoke(ctx, tt.query)

			tt.then(t, result, err)
			teamRepository.AssertExpectations(t)
		})
	}
}
