package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/account/usecases"
	"sportlink/api/domain/account"
	amocks "sportlink/mocks/api/domain/account"
)

func TestFindAccountUC_Invoke(t *testing.T) {
	testCases := []struct {
		name  string
		input usecases.FindAccountInput
		on    func(t *testing.T, repository *amocks.Repository, ctx context.Context)
		then  func(t *testing.T, result *[]account.Entity, err error)
	}{
		{
			name: "given repository returns accounts when finding by email then returns pointer to entities",
			input: usecases.FindAccountInput{
				Email: "user@example.com",
			},
			on: func(t *testing.T, repository *amocks.Repository, ctx context.Context) {
				entities := []account.Entity{
					{Email: "user@example.com", Nickname: "user1"},
				}
				repository.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q account.DomainQuery) bool {
						return len(q.Emails) == 1 && q.Emails[0] == "user@example.com" && len(q.Ids) == 0
					}),
				).Return(entities, nil)
			},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "user@example.com", (*result)[0].Email)
				assert.Equal(t, "user1", (*result)[0].Nickname)
			},
		},
		{
			name: "given repository returns accounts when finding by account id then returns pointer to entities",
			input: usecases.FindAccountInput{
				AccountID: "EMAIL#user@example.com",
			},
			on: func(t *testing.T, repository *amocks.Repository, ctx context.Context) {
				entities := []account.Entity{
					{ID: "EMAIL#user@example.com", Email: "user@example.com"},
				}
				repository.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q account.DomainQuery) bool {
						return len(q.Ids) == 1 && q.Ids[0] == "EMAIL#user@example.com" && len(q.Emails) == 0
					}),
				).Return(entities, nil)
			},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, *result, 1)
				assert.Equal(t, "EMAIL#user@example.com", (*result)[0].ID)
			},
		},
		{
			name: "given both account id and email when invoking then returns error without calling repository",
			input: usecases.FindAccountInput{
				AccountID: "EMAIL#a@b.com",
				Email:     "a@b.com",
			},
			on: func(t *testing.T, _ *amocks.Repository, _ context.Context) {},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "only one of account id or email must be provided")
			},
		},
		{
			name:  "given neither account id nor email when invoking then returns error without calling repository",
			input: usecases.FindAccountInput{},
			on: func(t *testing.T, _ *amocks.Repository, _ context.Context) {},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "account id or email is required")
			},
		},
		{
			name: "given whitespace-only fields when invoking then returns error without calling repository",
			input: usecases.FindAccountInput{
				AccountID: "   ",
				Email:     "\t",
			},
			on: func(t *testing.T, _ *amocks.Repository, _ context.Context) {},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "account id or email is required")
			},
		},
		{
			name: "given repository returns empty when finding then returns empty slice pointer",
			input: usecases.FindAccountInput{
				Email: "nobody@example.com",
			},
			on: func(t *testing.T, repository *amocks.Repository, ctx context.Context) {
				repository.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q account.DomainQuery) bool {
						return len(q.Emails) == 1 && q.Emails[0] == "nobody@example.com"
					}),
				).Return([]account.Entity{}, nil)
			},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Empty(t, *result)
			},
		},
		{
			name: "given repository find fails when finding then returns error",
			input: usecases.FindAccountInput{
				AccountID: "EMAIL#missing@example.com",
			},
			on: func(t *testing.T, repository *amocks.Repository, ctx context.Context) {
				repository.On("Find",
					mock.MatchedBy(func(c context.Context) bool { return c == ctx }),
					mock.MatchedBy(func(q account.DomainQuery) bool {
						return len(q.Ids) == 1 && q.Ids[0] == "EMAIL#missing@example.com"
					}),
				).Return(nil, errors.New("database unavailable"))
			},
			then: func(t *testing.T, result *[]account.Entity, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database unavailable", err.Error())
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up
			ctx := context.Background()
			repository := amocks.NewRepository(t)
			uc := usecases.NewFindAccountUC(repository)

			// given
			tt.on(t, repository, ctx)

			// when
			result, err := uc.Invoke(ctx, tt.input)

			// then
			tt.then(t, result, err)
		})
	}
}
