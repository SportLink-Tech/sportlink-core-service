package account_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/account/usecases"
	"sportlink/api/domain/account"
	"sportlink/api/infrastructure/middleware"
	caccount "sportlink/api/infrastructure/rest/account"
	amocks "sportlink/mocks/api/application"
)

func TestFindAccount(t *testing.T) {
	testCases := []struct {
		name string
		path string
		on   func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context)
		then func(t *testing.T, code int, body interface{})
	}{
		{
			name: "given account exists when getting by account id then returns account details",
			path: "/account?account_id=EMAIL%23user%40example.com",
			on: func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context) {
				slice := []account.Entity{
					{
						ID:       "EMAIL#user@example.com",
						Email:    "user@example.com",
						Nickname: "user1",
					},
				}
				uc.On("Invoke",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(in usecases.FindAccountInput) bool {
						return in.AccountID == "EMAIL#user@example.com" && in.Email == ""
					}),
				).Return(&slice, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
				first := arr[0].(map[string]interface{})
				assert.Equal(t, "EMAIL#user@example.com", first["ID"])
				assert.Equal(t, "user@example.com", first["Email"])
				assert.Equal(t, "user1", first["Nickname"])
			},
		},
		{
			name: "given account exists when getting by email then returns account details",
			path: "/account?email=user%40example.com",
			on: func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context) {
				slice := []account.Entity{
					{ID: "EMAIL#user@example.com", Email: "user@example.com", Nickname: "user1"},
				}
				uc.On("Invoke",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(in usecases.FindAccountInput) bool {
						return in.AccountID == "" && in.Email == "user@example.com"
					}),
				).Return(&slice, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusOK, code)
				arr := body.([]interface{})
				assert.Len(t, arr, 1)
				first := arr[0].(map[string]interface{})
				assert.Equal(t, "user@example.com", first["Email"])
			},
		},
		{
			name: "given no account id or email query when getting then returns validation error",
			path: "/account",
			on:   func(t *testing.T, _ *amocks.UseCase[usecases.FindAccountInput, []account.Entity], _ context.Context) {},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
				assert.Contains(t, m["message"], "account_id or email")
			},
		},
		{
			name: "given only whitespace query params when getting then returns validation error",
			path: "/account?account_id=%20&email=%09",
			on:   func(t *testing.T, _ *amocks.UseCase[usecases.FindAccountInput, []account.Entity], _ context.Context) {},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
			},
		},
		{
			name: "given both account id and email query when getting then returns validation error",
			path: "/account?account_id=id&email=e%40e.com",
			on:   func(t *testing.T, _ *amocks.UseCase[usecases.FindAccountInput, []account.Entity], _ context.Context) {},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusBadRequest, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", m["code"])
				assert.Contains(t, m["message"], "only one of account_id or email")
			},
		},
		{
			name: "given no account when getting then returns not found",
			path: "/account?account_id=EMAIL%23missing%40example.com",
			on: func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context) {
				empty := []account.Entity{}
				uc.On("Invoke",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(in usecases.FindAccountInput) bool {
						return in.AccountID == "EMAIL#missing@example.com" && in.Email == ""
					}),
				).Return(&empty, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusNotFound, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "not_found", m["code"])
				assert.Equal(t, "No account found", m["message"])
			},
		},
		{
			name: "given lookup returns no result payload when getting then returns not found",
			path: "/account?account_id=any-id",
			on: func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context) {
				uc.On("Invoke",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(in usecases.FindAccountInput) bool {
						return in.AccountID == "any-id" && in.Email == ""
					}),
				).Return(nil, nil)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusNotFound, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "not_found", m["code"])
			},
		},
		{
			name: "given lookup fails when getting then returns conflict",
			path: "/account?account_id=EMAIL%23broken%40example.com",
			on: func(t *testing.T, uc *amocks.UseCase[usecases.FindAccountInput, []account.Entity], ctx context.Context) {
				uc.On("Invoke",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(in usecases.FindAccountInput) bool {
						return in.AccountID == "EMAIL#broken@example.com" && in.Email == ""
					}),
				).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, code int, body interface{}) {
				assert.Equal(t, http.StatusConflict, code)
				m := body.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", m["code"])
				assert.Contains(t, m["message"], "use case execution failed")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// set up
			ucMock := amocks.NewUseCase[usecases.FindAccountInput, []account.Entity](t)
			ctl := caccount.NewController(ucMock)

			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.Use(middleware.ErrorHandler())
			r.GET("/account", ctl.Find)

			// given
			tc.on(t, ucMock, ctx)

			// when
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			var body interface{}
			if rec.Body.Len() > 0 {
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
			}

			// then
			tc.then(t, rec.Code, body)
		})
	}
}
