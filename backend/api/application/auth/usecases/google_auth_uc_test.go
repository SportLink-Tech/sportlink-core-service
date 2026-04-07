package usecases_test

import (
	"context"
	"fmt"
	"sportlink/api/application/auth/service"
	"sportlink/api/application/auth/usecases"
	"sportlink/api/domain/account"
	amocks "sportlink/mocks/api/domain/account"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- mocks ---

type mockGoogleVerifier struct{ mock.Mock }

func (m *mockGoogleVerifier) Verify(ctx context.Context, idToken string) (*service.GoogleTokenInfo, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.GoogleTokenInfo), args.Error(1)
}

type mockJWTService struct{ mock.Mock }

func (m *mockJWTService) Generate(accountID string) (string, error) {
	args := m.Called(accountID)
	return args.String(0), args.Error(1)
}

// --- tests ---

var validTokenInfo = &service.GoogleTokenInfo{
	Sub:           "google-sub-123",
	Email:         "user@gmail.com",
	EmailVerified: "true",
	Name:          "Jorge Cabrera",
	Picture:       "https://photo.url",
}

func TestGoogleAuthUC_Invoke(t *testing.T) {
	tests := []struct {
		name    string
		idToken string
		on      func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService)
		then    func(t *testing.T, result *usecases.GoogleAuthResult, err error)
	}{
		{
			name:    "new user: creates account and returns token",
			idToken: "valid-id-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-id-token").Return(validTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.MatchedBy(func(q account.DomainQuery) bool {
					return len(q.Emails) == 1 && q.Emails[0] == "user@gmail.com"
				})).Return([]account.Entity{}, nil)
				repo.On("Save", mock.Anything, mock.MatchedBy(func(e account.Entity) bool {
					return e.Email == "user@gmail.com" && e.Picture == "https://photo.url" && e.AccountID != ""
				})).Return(nil)
				jwt.On("Generate", mock.AnythingOfType("string")).Return("signed-jwt", nil)
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "signed-jwt", result.JWTToken)
				assert.NotEmpty(t, result.AccountID)
			},
		},
		{
			name:    "existing user: returns token without creating account",
			idToken: "valid-id-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-id-token").Return(validTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{
					{ID: "EMAIL#user@gmail.com", AccountID: "01JQTEST0000000000000000AB", Email: "user@gmail.com"},
				}, nil)
				jwt.On("Generate", "01JQTEST0000000000000000AB").Return("signed-jwt", nil)
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "signed-jwt", result.JWTToken)
				assert.Equal(t, "01JQTEST0000000000000000AB", result.AccountID)
			},
		},
		{
			name:    "fails when Google token is invalid",
			idToken: "bad-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "bad-token").Return(nil, fmt.Errorf("invalid token"))
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "invalid Google token")
			},
		},
		{
			name:    "fails when account repository returns error",
			idToken: "valid-id-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-id-token").Return(validTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{}, fmt.Errorf("db error"))
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error finding account")
			},
		},
		{
			name:    "fails when saving new account returns error",
			idToken: "valid-id-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-id-token").Return(validTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{}, nil)
				repo.On("Save", mock.Anything, mock.Anything).Return(fmt.Errorf("save error"))
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error creating account")
			},
		},
		{
			name:    "fails when JWT generation returns error",
			idToken: "valid-id-token",
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-id-token").Return(validTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{
					{ID: "EMAIL#user@gmail.com", AccountID: "01JQTEST0000000000000000AB", Email: "user@gmail.com"},
				}, nil)
				jwt.On("Generate", "01JQTEST0000000000000000AB").Return("", fmt.Errorf("jwt error"))
			},
			then: func(t *testing.T, result *usecases.GoogleAuthResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "error generating token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			verifier := &mockGoogleVerifier{}
			repo := amocks.NewRepository(t)
			jwtSvc := &mockJWTService{}

			uc := usecases.NewGoogleAuthUC(verifier, repo, jwtSvc)
			tt.on(verifier, repo, jwtSvc)

			result, err := uc.Invoke(context.Background(), tt.idToken)

			tt.then(t, result, err)
		})
	}
}
