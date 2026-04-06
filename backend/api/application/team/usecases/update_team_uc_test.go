package usecases_test

import (
	"context"
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

func TestUpdateTeamUC_Invoke(t *testing.T) {
	existingTeam := team.Entity{
		ID:       "SPORT#Football#NAME#Boca Juniors",
		Name:     "Boca Juniors",
		Sport:    common.Football,
		Category: common.L1,
		Stats:    *common.NewStats(0, 0, 0),
		Members:  []player.Entity{},
	}

	tests := []struct {
		name  string
		input team.PatchInput
		on    func(t *testing.T, repo *mmocks.Repository)
		then  func(t *testing.T, result *team.Entity, err error)
	}{
		{
			name: "updates team name successfully",
			input: team.PatchInput{
				ID:   team.ID{Sport: common.Football, Name: "Boca Juniors"},
				Name: strPtr("Boca Senior"),
			},
			on: func(t *testing.T, repo *mmocks.Repository) {
				repo.On("Find", mock.Anything, mock.MatchedBy(func(q team.DomainQuery) bool {
					return q.Name == "Boca Juniors" && len(q.Sports) == 1 && q.Sports[0] == common.Football
				})).Return([]team.Entity{existingTeam}, nil)

				repo.On("Update", mock.Anything, "SPORT#Football#NAME#Boca Juniors", mock.MatchedBy(func(e team.Entity) bool {
					return e.Name == "Boca Senior" && e.ID == "SPORT#Football#NAME#Boca Senior"
				})).Return(nil)
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Boca Senior", result.Name)
				assert.Equal(t, "SPORT#Football#NAME#Boca Senior", result.ID)
			},
		},
		{
			name: "does not change name when not provided in patch",
			input: team.PatchInput{
				ID: team.ID{Sport: common.Football, Name: "Boca Juniors"},
			},
			on: func(t *testing.T, repo *mmocks.Repository) {
				repo.On("Find", mock.Anything, mock.Anything).Return([]team.Entity{existingTeam}, nil)
				repo.On("Update", mock.Anything, "SPORT#Football#NAME#Boca Juniors", mock.MatchedBy(func(e team.Entity) bool {
					return e.Name == "Boca Juniors"
				})).Return(nil)
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Boca Juniors", result.Name)
			},
		},
		{
			name: "returns error when team not found",
			input: team.PatchInput{
				ID:   team.ID{Sport: common.Football, Name: "Unknown"},
				Name: strPtr("New Name"),
			},
			on: func(t *testing.T, repo *mmocks.Repository) {
				repo.On("Find", mock.Anything, mock.Anything).Return([]team.Entity{}, nil)
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "team not found")
				assert.Nil(t, result)
			},
		},
		{
			name: "returns error when find fails",
			input: team.PatchInput{
				ID: team.ID{Sport: common.Football, Name: "Boca Juniors"},
			},
			on: func(t *testing.T, repo *mmocks.Repository) {
				repo.On("Find", mock.Anything, mock.Anything).Return([]team.Entity{}, fmt.Errorf("db error"))
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error finding team")
				assert.Nil(t, result)
			},
		},
		{
			name: "returns error when update fails",
			input: team.PatchInput{
				ID:   team.ID{Sport: common.Football, Name: "Boca Juniors"},
				Name: strPtr("Boca Senior"),
			},
			on: func(t *testing.T, repo *mmocks.Repository) {
				repo.On("Find", mock.Anything, mock.Anything).Return([]team.Entity{existingTeam}, nil)
				repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("update failed"))
			},
			then: func(t *testing.T, result *team.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error updating team")
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mmocks.Repository{}
			uc := usecases.NewUpdateTeamUC(repo)

			tt.on(t, repo)

			result, err := uc.Invoke(context.Background(), tt.input)

			tt.then(t, result, err)
		})
	}
}

func strPtr(s string) *string { return &s }
