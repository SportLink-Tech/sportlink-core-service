package account_test

import (
	"context"
	daccount "sportlink/api/domain/account"
	"sportlink/api/infrastructure/persistence/account"
	"sportlink/dev/testcontainer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Save(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := account.NewRepository(dynamoDbClient, "SportLinkCore")

	testCases := []struct {
		name       string
		entity     daccount.Entity
		assertions func(t *testing.T, err error)
	}{
		{
			name:   "save account successfully",
			entity: daccount.NewAccount("test@example.com", "testuser"),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name:   "save account with special characters in email successfully",
			entity: daccount.NewAccount("test.user+tag@example.co.uk", "testuser2"),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name:   "save account with nickname containing spaces successfully",
			entity: daccount.NewAccount("user@example.com", "test user"),
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := repository.Save(ctx, testCase.entity)

			// then
			testCase.assertions(t, err)
		})
	}
}

func Test_Find(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := account.NewRepository(dynamoDbClient, "SportLinkCore")

	testCases := []struct {
		name       string
		query      daccount.DomainQuery
		on         func(t *testing.T, repository daccount.Repository)
		assertions func(t *testing.T, entities []daccount.Entity, err error)
	}{
		{
			name: "find account by email successfully",
			query: daccount.DomainQuery{
				Emails: []string{"test@example.com"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				if err := repository.Save(ctx, daccount.NewAccount("test@example.com", "testuser")); err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name: "find account by multiple emails successfully",
			query: daccount.DomainQuery{
				Emails: []string{"user1@example.com", "user2@example.com"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				for _, acc := range []daccount.Entity{
					daccount.NewAccount("user1@example.com", "user1"),
					daccount.NewAccount("user2@example.com", "user2"),
					daccount.NewAccount("user3@example.com", "user3"),
				} {
					if err := repository.Save(ctx, acc); err != nil {
						t.Fatal(err)
					}
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				emails := []string{entities[0].Email, entities[1].Email}
				assert.Contains(t, emails, "user1@example.com")
				assert.Contains(t, emails, "user2@example.com")
				assert.NotContains(t, emails, "user3@example.com")
			},
		},
		{
			name: "find account by id successfully",
			query: daccount.DomainQuery{
				Ids: []string{"EMAIL#test@example.com"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				if err := repository.Save(ctx, daccount.NewAccount("test@example.com", "testuser")); err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "EMAIL#test@example.com", entities[0].ID)
				assert.Equal(t, "test@example.com", entities[0].Email)
			},
		},
		{
			name: "find account by multiple ids successfully",
			query: daccount.DomainQuery{
				Ids: []string{"EMAIL#user1@example.com", "EMAIL#user2@example.com"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				for _, acc := range []daccount.Entity{
					daccount.NewAccount("user1@example.com", "user1"),
					daccount.NewAccount("user2@example.com", "user2"),
				} {
					if err := repository.Save(ctx, acc); err != nil {
						t.Fatal(err)
					}
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				ids := []string{entities[0].ID, entities[1].ID}
				assert.Contains(t, ids, "EMAIL#user1@example.com")
				assert.Contains(t, ids, "EMAIL#user2@example.com")
			},
		},
		{
			name: "find account by email and nickname successfully",
			query: daccount.DomainQuery{
				Emails:    []string{"test@example.com"},
				Nicknames: []string{"testuser"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				for _, acc := range []daccount.Entity{
					daccount.NewAccount("test@example.com", "testuser"),
					daccount.NewAccount("other@example.com", "otheruser"),
				} {
					if err := repository.Save(ctx, acc); err != nil {
						t.Fatal(err)
					}
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name: "find account by email and multiple nicknames successfully",
			query: daccount.DomainQuery{
				Emails:    []string{"test@example.com"},
				Nicknames: []string{"testuser", "anotheruser"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				if err := repository.Save(ctx, daccount.NewAccount("test@example.com", "testuser")); err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name: "find account that does not exist",
			query: daccount.DomainQuery{
				Emails: []string{"nonexistent@example.com"},
			},
			on: func(t *testing.T, repository daccount.Repository) {},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name:  "find with no criteria returns empty",
			query: daccount.DomainQuery{},
			on: func(t *testing.T, repository daccount.Repository) {
				if err := repository.Save(ctx, daccount.NewAccount("test@example.com", "testuser")); err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "find with only nickname returns empty",
			query: daccount.DomainQuery{
				Nicknames: []string{"testuser"},
			},
			on: func(t *testing.T, repository daccount.Repository) {
				if err := repository.Save(ctx, daccount.NewAccount("test@example.com", "testuser")); err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.on(t, repository)
			entities, err := repository.Find(ctx, testCase.query)
			testCase.assertions(t, entities, err)
		})
	}
}

func Test_FindByAccountID(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := account.NewRepository(dynamoDbClient, "SportLinkCore")

	acc := daccount.NewAccount("ulid@example.com", "uliduser")
	if err := repository.Save(ctx, acc); err != nil {
		t.Fatal(err)
	}

	entities, err := repository.Find(ctx, daccount.DomainQuery{
		AccountIDs: []string{acc.AccountID},
	})

	assert.Nil(t, err)
	assert.Len(t, entities, 1)
	assert.Equal(t, acc.AccountID, entities[0].AccountID)
	assert.Equal(t, "ulid@example.com", entities[0].Email)
}
