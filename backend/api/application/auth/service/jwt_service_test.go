package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sportlink/api/application/auth/service"
)

func TestJwtService_Generate(t *testing.T) {

	tests := []struct {
		name      string
		secret    string
		accountID string
		then      func(t *testing.T, token string, err error)
	}{
		{
			name:      "given valid account id when generating token then returns no error",
			secret:    "secret",
			accountID: "account-id",
			then: func(t *testing.T, token string, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// set up
			svc := service.NewJWTService(tt.secret)

			// when
			token, err := svc.Generate(tt.accountID)

			// then
			tt.then(t, token, err)
		})
	}
}

func TestJwtService_Parse(t *testing.T) {

	tests := []struct {
		name        string
		parseSecret string
		signSecret  string
		accountID   string
		fixedToken  string
		then        func(t *testing.T, claims *service.AccessTokenClaims, err error)
	}{
		{
			name:        "given valid token for an account when parsing then returns that account identifier",
			parseSecret: "secret",
			signSecret:  "secret",
			accountID:   "account-id",
			then: func(t *testing.T, claims *service.AccessTokenClaims, err error) {
				t.Helper()
				assert.NoError(t, err)
				require.NotNil(t, claims)
				assert.Equal(t, "account-id", claims.AccountID)
				require.NotNil(t, claims.ExpiresAt)
				require.NotNil(t, claims.IssuedAt)
				assert.Greater(t, claims.ExpiresAt.Unix(), claims.IssuedAt.Unix())
			},
		},
		{
			name:        "given token and different verification secret when parsing then returns error",
			parseSecret: "other-secret",
			signSecret:  "secret",
			accountID:   "account-id",
			then: func(t *testing.T, claims *service.AccessTokenClaims, err error) {
				t.Helper()
				assert.Error(t, err)
			},
		},
		{
			name:        "given malformed token when parsing then returns error",
			parseSecret: "secret",
			fixedToken:  "not-a-jwt",
			then: func(t *testing.T, claims *service.AccessTokenClaims, err error) {
				t.Helper()
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// set up
			var token string
			if tt.fixedToken != "" {
				token = tt.fixedToken
			} else {
				signer := service.NewJWTService(tt.signSecret)
				var err error
				token, err = signer.Generate(tt.accountID)
				assert.NoError(t, err)
			}
			svc := service.NewJWTService(tt.parseSecret)

			// when
			claims, err := svc.Parse(token)

			// then
			tt.then(t, claims, err)
		})
	}
}
