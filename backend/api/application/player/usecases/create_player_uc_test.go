package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/player/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	pmocks "sportlink/mocks/api/domain/player"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreatePlayerUC(t *testing.T) {
	ctx := context.Background()
	saveErr := errors.New("database error")

	tests := []struct {
		name  string
		input player.Entity
		given func(t *testing.T, repository *pmocks.Repository)
		then  func(t *testing.T, result *player.Entity, err error)
	}{
		{
			name:  "given save succeeds when invoke then returns entity with generated id",
			input: player.NewPlayer(common.L1, common.Football),
			given: func(t *testing.T, repository *pmocks.Repository) {
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity player.Entity) bool {
					return entity.ID != "" &&
						entity.Category == common.L1 &&
						entity.Sport == common.Football
				})).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.ID)
				assert.EqualValues(t, common.L1, result.Category)
				assert.Equal(t, common.Football, result.Sport)
			},
		},
		{
			name:  "given repository save fails when invoke then returns wrapped error",
			input: player.NewPlayer(common.L3, common.Tennis),
			given: func(t *testing.T, repository *pmocks.Repository) {
				repository.On("Save", mock.Anything, mock.Anything).Return(saveErr)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting player in database")
				assert.ErrorIs(t, err, saveErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repository := &pmocks.Repository{}
			uc := usecases.NewCreatePlayerUC(repository)

			tt.given(t, repository)

			result, err := uc.Invoke(ctx, tt.input)

			tt.then(t, result, err)
			repository.AssertExpectations(t)
		})
	}
}
