package usecases_test

import (
	"context"
	"errors"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchrequest"
	mrmocks "sportlink/mocks/api/domain/matchrequest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewUpdateMatchRequestStatusUC(t *testing.T) {
	ctx := context.Background()
	updateErr := errors.New("conditional check failed")

	tests := []struct {
		name  string
		input usecases.UpdateMatchRequestStatusInput
		given func(t *testing.T, mrRepository *mrmocks.Repository)
		then  func(t *testing.T, err error)
	}{
		{
			name: "given update succeeds when invoke then returns nil",
			input: usecases.UpdateMatchRequestStatusInput{
				ID: "req-1", OwnerAccountID: "owner-1", NewStatus: matchrequest.StatusAccepted,
			},
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("UpdateStatus", mock.Anything, "req-1", "owner-1", matchrequest.StatusAccepted).Return(nil)
			},
			then: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "given repository fails when invoke then returns wrapped error",
			input: usecases.UpdateMatchRequestStatusInput{
				ID: "req-1", OwnerAccountID: "owner-1", NewStatus: matchrequest.StatusRejected,
			},
			given: func(t *testing.T, mrRepository *mrmocks.Repository) {
				mrRepository.On("UpdateStatus", mock.Anything, "req-1", "owner-1", matchrequest.StatusRejected).Return(updateErr)
			},
			then: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "error while updating match request status")
				assert.ErrorIs(t, err, updateErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mrRepository := &mrmocks.Repository{}
			uc := usecases.NewUpdateMatchRequestStatusUC(mrRepository)

			tt.given(t, mrRepository)

			err := uc.Invoke(ctx, tt.input)

			tt.then(t, err)
			mrRepository.AssertExpectations(t)
		})
	}
}
