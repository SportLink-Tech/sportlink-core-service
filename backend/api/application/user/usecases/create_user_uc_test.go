package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/user/usecases"
	"sportlink/api/domain/user"
	umocks "sportlink/mocks/api/domain/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreateUserUC(t *testing.T) {
	ctx := context.Background()
	saveErr := errors.New("database write failed")

	tests := []struct {
		name  string
		input user.Entity
		given func(t *testing.T, repository *umocks.Repository)
		then  func(t *testing.T, result *user.Entity, err error)
	}{
		{
			name:  "given save succeeds when invoke then returns same entity pointer",
			input: user.NewUser("Jane", "Doe", []string{"player-1"}),
			given: func(t *testing.T, repository *umocks.Repository) {
				repository.On("Save", mock.Anything, mock.MatchedBy(func(e user.Entity) bool {
					return e.FirstName == "Jane" && e.LastName == "Doe" &&
						len(e.PlayerIDs) == 1 && e.PlayerIDs[0] == "player-1" && e.ID != ""
				})).Return(nil)
			},
			then: func(t *testing.T, result *user.Entity, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "Jane", result.FirstName)
				assert.Equal(t, "Doe", result.LastName)
				assert.NotEmpty(t, result.ID)
			},
		},
		{
			name:  "given repository save fails when invoke then returns wrapped error",
			input: user.NewUser("John", "Smith", nil),
			given: func(t *testing.T, repository *umocks.Repository) {
				repository.On("Save", mock.Anything, mock.MatchedBy(func(e user.Entity) bool {
					return e.FirstName == "John" && e.LastName == "Smith"
				})).Return(saveErr)
			},
			then: func(t *testing.T, result *user.Entity, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				assert.Contains(t, err.Error(), "error while inserting user in database")
				assert.ErrorIs(t, err, saveErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repository := &umocks.Repository{}
			uc := usecases.NewCreateUserUC(repository)

			tt.given(t, repository)

			result, err := uc.Invoke(ctx, tt.input)

			tt.then(t, result, err)
			repository.AssertExpectations(t)
		})
	}
}
