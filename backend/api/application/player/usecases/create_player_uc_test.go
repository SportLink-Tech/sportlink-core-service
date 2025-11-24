package usecases_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sportlink/api/application/player/usecases"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	mocks "sportlink/mocks/api/domain/player"
	"testing"
)

func TestCreatePlayerUC_Invoke(t *testing.T) {

	tests := []struct {
		name  string
		input player.Entity
		on    func(t *testing.T, repository *mocks.Repository)
		then  func(t *testing.T, result *player.Entity, err error)
	}{
		{
			name: "save player successfully",
			input: player.Entity{
				ID:       "player1",
				Category: common.L1,
				Sport:    common.Football,
			},
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Find", mock.MatchedBy(func(query player.DomainQuery) bool {
					return query.Id == "player1" && query.Category == common.L1 && query.Sport == common.Football
				})).Return([]player.Entity{}, fmt.Errorf("not found"))
				repository.On("Save", mock.MatchedBy(func(entity player.Entity) bool {
					return entity.ID == "player1" && entity.Category == common.L1 && entity.Sport == common.Football
				})).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "player1", result.ID)
				assert.EqualValues(t, common.L1, result.Category)
				assert.Equal(t, common.Football, result.Sport)
			},
		},
		{
			name: "player already exists returns the existing player",
			input: player.Entity{
				ID:       "player2",
				Category: common.L2,
				Sport:    common.Paddle,
			},
			on: func(t *testing.T, repository *mocks.Repository) {
				existingPlayer := player.Entity{
					ID:       "player2",
					Category: common.L2,
					Sport:    common.Paddle,
				}
				repository.On("Find", mock.MatchedBy(func(query player.DomainQuery) bool {
					return query.Id == "player2"
				})).Return([]player.Entity{existingPlayer}, nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "Player already exist")
			},
		},
		{
			name: "error while saving player returns error",
			input: player.Entity{
				ID:       "player3",
				Category: common.L3,
				Sport:    common.Tennis,
			},
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Find", mock.Anything).Return([]player.Entity{}, fmt.Errorf("not found"))
				repository.On("Save", mock.Anything).Return(fmt.Errorf("database error"))
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting player in database")
			},
		},
		{
			name: "save paddle player successfully",
			input: player.Entity{
				ID:       "paddle_player1",
				Category: common.L4,
				Sport:    common.Paddle,
			},
			on: func(t *testing.T, repository *mocks.Repository) {
				repository.On("Find", mock.Anything).Return([]player.Entity{}, fmt.Errorf("not found"))
				repository.On("Save", mock.Anything).Return(nil)
			},
			then: func(t *testing.T, result *player.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "paddle_player1", result.ID)
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
			result, err := uc.Invoke(tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
