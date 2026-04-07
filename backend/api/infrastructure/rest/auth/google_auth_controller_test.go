package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sportlink/api/application/auth/service"
	"sportlink/api/application/auth/usecases"
	"sportlink/api/domain/account"
	"sportlink/api/infrastructure/middleware"
	cauth "sportlink/api/infrastructure/rest/auth"
	amocks "sportlink/mocks/api/domain/account"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mocks

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

var googleTokenInfo = &service.GoogleTokenInfo{
	Email: "user@gmail.com", EmailVerified: "true",
	Name: "Jorge Cabrera", Picture: "https://photo.url",
}

func TestGoogleAuthController(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name       string
		body       map[string]interface{}
		on         func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService)
		assertions func(t *testing.T, code int, body map[string]interface{}, resp *httptest.ResponseRecorder)
	}{
		{
			name: "new user authenticates and receives cookie",
			body: map[string]interface{}{"id_token": "valid-token"},
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-token").Return(googleTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{}, nil)
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				jwt.On("Generate", mock.Anything).Return("jwt-token", nil)
			},
			assertions: func(t *testing.T, code int, body map[string]interface{}, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, code)
				assert.NotEmpty(t, body["account_id"])
				assert.Contains(t, resp.Header().Get("Set-Cookie"), "token=jwt-token")
				assert.Contains(t, resp.Header().Get("Set-Cookie"), "HttpOnly")
			},
		},
		{
			name: "existing user authenticates and receives cookie",
			body: map[string]interface{}{"id_token": "valid-token"},
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "valid-token").Return(googleTokenInfo, nil)
				repo.On("Find", mock.Anything, mock.Anything).Return([]account.Entity{
					{ID: "EMAIL#user@gmail.com", AccountID: "01JQTEST0000000000000000AB", Email: "user@gmail.com"},
				}, nil)
				jwt.On("Generate", "01JQTEST0000000000000000AB").Return("jwt-token", nil)
			},
			assertions: func(t *testing.T, code int, body map[string]interface{}, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, code)
				assert.Equal(t, "01JQTEST0000000000000000AB", body["account_id"])
			},
		},
		{
			name:       "fails with missing id_token",
			body:       map[string]interface{}{},
			on:         func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {},
			assertions: func(t *testing.T, code int, body map[string]interface{}, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, code)
			},
		},
		{
			name: "fails with invalid Google token",
			body: map[string]interface{}{"id_token": "bad-token"},
			on: func(verifier *mockGoogleVerifier, repo *amocks.Repository, jwt *mockJWTService) {
				verifier.On("Verify", mock.Anything, "bad-token").Return(nil, assert.AnError)
			},
			assertions: func(t *testing.T, code int, body map[string]interface{}, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, code)
				assert.Equal(t, "unauthorized", body["code"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			verifier := &mockGoogleVerifier{}
			repo := amocks.NewRepository(t)
			jwtSvc := &mockJWTService{}

			uc := usecases.NewGoogleAuthUC(verifier, repo, jwtSvc)
			controller := cauth.NewController(uc, v)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.POST("/auth/google", controller.GoogleAuth)

			tc.on(verifier, repo, jwtSvc)
			jsonData, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/auth/google", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			var body map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &body)
			tc.assertions(t, resp.Code, body, resp)
		})
	}
}
