package usecases_test

import (
	"context"
	"fmt"
	"sportlink/api/application/player/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	mocks "sportlink/mocks/api/domain/player"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlayerUC_Invoke(t *testing.T) {

	tests := []struct {
		name  string
		input player.Entity
		on    func(t *testing.T, repository *mocks.Repository)
		then  func(t *testing.T, result *player.Entity, err error)
	}{
		{
			name:  "save player successfully",
			input: player.NewPlayer(common.L1, common.Football),
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity player.Entity) bool {
					return entity.ID != "" && // ULID is generated
						entity.Category == common.L1 &&
						entity.Sport == common.Football
				})).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID) // ULID is generated
				assert.EqualValues(t, common.L1, result.Category)
				assert.Equal(t, common.Football, result.Sport)
			},
		},
		{
			name:  "save player with ULID generated automatically",
			input: player.NewPlayer(common.L2, common.Paddle),
			on: func(t *testing.T, repository *mocks.Repository) {
				// With ULID, each player gets a unique ID, so duplicates by ID are not possible
				repository.On("Save", mock.Anything, mock.MatchedBy(func(entity player.Entity) bool {
					return entity.ID != "" && // ULID is generated
						entity.Category == common.L2 &&
						entity.Sport == common.Paddle
				})).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID) // ULID is generated
				assert.EqualValues(t, common.L2, result.Category)
				assert.Equal(t, common.Paddle, result.Sport)
			},
		},
		{
			name:  "error while saving player returns error",
			input: player.NewPlayer(common.L3, common.Tennis),
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Save", mock.Anything, mock.Anything).Return(fmt.Errorf("database error"))
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting player in database")
			},
		},
		{
			name:  "save paddle player successfully",
			input: player.NewPlayer(common.L4, common.Paddle),
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID) // ULID is generated
				assert.Equal(t, common.Paddle, result.Sport)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			repository := &mocks.Repository{}
			uc := usecases.NewCreatePlayerUC(repository)

			// given
			tt.on(t, repository)

			// when
			result, err := uc.Invoke(context.Background(), tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
