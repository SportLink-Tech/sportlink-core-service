package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchrequest/usecases"
	domainreq "sportlink/api/domain/matchrequest"
	reqmocks "sportlink/mocks/api/domain/matchrequest"
)

func TestUpdateMatchRequestStatusUC_Invoke(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name  string
		input usecases.UpdateMatchRequestStatusInput
		on    func(t *testing.T, repository *reqmocks.Repository)
		then  func(t *testing.T, err error)
	}{
		{
			name: "given repository succeeds when updating status then returns no error",
			input: usecases.UpdateMatchRequestStatusInput{
				ID:             "mr-1",
				OwnerAccountID: "owner-acc",
				NewStatus:      domainreq.StatusAccepted,
			},
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("UpdateStatus", mock.Anything, "mr-1", "owner-acc", domainreq.StatusAccepted).Return(nil)
			},
			then: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "given repository fails when updating status then returns wrapped error",
			input: usecases.UpdateMatchRequestStatusInput{
				ID:             "mr-1",
				OwnerAccountID: "owner-acc",
				NewStatus:      domainreq.StatusRejected,
			},
			on: func(t *testing.T, repository *reqmocks.Repository) {
				repository.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("conditional check failed"))
			},
			then: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error while updating match request status")
				assert.Contains(t, err.Error(), "conditional check failed")
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			repository := reqmocks.NewRepository(t)
			uc := usecases.NewUpdateMatchRequestStatusUC(repository)

			// given
			tt.on(t, repository)

			// when
			err := uc.Invoke(ctx, tt.input)

			// then
			tt.then(t, err)
		})
	}
}
