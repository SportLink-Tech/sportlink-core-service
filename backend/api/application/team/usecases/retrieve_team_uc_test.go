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

func TestNewRetrieveTeamUC(t *testing.T) {
	ctx := context.Background()
	findErr := errors.New("database error")

	tests := []struct {
		name  string
		id    team.ID
		given func(t *testing.T, teamRepository *mmocks.Repository)
		then  func(t *testing.T, result *team.Entity, err error)
	}{
		{
			name: "given team exists when invoke then returns first entity",
			id:   team.ID{Name: "Boca Jr", Sport: common.Football},
			given: func(t *testing.T, teamRepository *mmocks.Repository) {
				found := team.Entity{
					Name: "Boca Jr", Sport: common.Football, Category: common.L1,
					Stats: *common.NewStats(10, 0, 0), Members: []player.Entity{},
				}
				teamRepository.On("Find", mock.Anything, mock.MatchedBy(func(q team.DomainQuery) bool {
					return q.Name == "Boca Jr" && len(q.Sports) == 1 && q.Sports[0] == common.Football
				})).Return([]team.Entity{found}, nil)
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "Boca Jr", result.Name)
				assert.Equal(t, common.Football, result.Sport)
			},
		},
		{
			name: "given find fails when invoke then returns error",
			id:   team.ID{Name: "Missing", Sport: common.Paddle},
			given: func(t *testing.T, teamRepository *mmocks.Repository) {
				teamRepository.On("Find", mock.Anything, mock.Anything).Return([]team.Entity{}, findErr)
			},
			then: func(t *testing.T, result *team.Entity, err error) {
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
			uc := usecases.NewRetrieveTeamUC(teamRepository)

			tt.given(t, teamRepository)

			result, err := uc.Invoke(ctx, tt.id)

			tt.then(t, result, err)
			teamRepository.AssertExpectations(t)
		})
	}
}
