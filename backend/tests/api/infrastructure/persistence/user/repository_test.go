package user_test

import (
	"context"
	duser "sportlink/api/domain/user"
	"sportlink/api/infrastructure/persistence/user"
	"sportlink/dev/testcontainer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Save(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := user.NewRepository(dynamoDbClient, "SportLinkCore")

	testCases := []struct {
		name       string
		entity     duser.Entity
		assertions func(t *testing.T, err error)
	}{
		{
			name: "save user successfully",
			entity: duser.Entity{
				ID:        "user123",
				FirstName: "John",
				LastName:  "Doe",
				PlayerIDs: []string{"player1", "player2"},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save user with empty PlayerIDs successfully",
			entity: duser.Entity{
				ID:        "user456",
				FirstName: "Jane",
				LastName:  "Smith",
				PlayerIDs: []string{},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save user with single PlayerID successfully",
			entity: duser.Entity{
				ID:        "user789",
				FirstName: "Bob",
				LastName:  "Johnson",
				PlayerIDs: []string{"player3"},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save user with special characters in ID successfully",
			entity: duser.Entity{
				ID:        "user-123_test",
				FirstName: "Alice",
				LastName:  "Williams",
				PlayerIDs: []string{"player4"},
			},
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
	repository := user.NewRepository(dynamoDbClient, "SportLinkCore")

	testCases := []struct {
		name       string
		query      duser.DomainQuery
		on         func(t *testing.T, repository duser.Repository)
		assertions func(t *testing.T, entities []duser.Entity, err error)
	}{
		{
			name: "find user by id successfully",
			query: duser.DomainQuery{
				Ids: []string{"user123"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
				assert.Equal(t, "John", entities[0].FirstName)
				assert.Equal(t, "Doe", entities[0].LastName)
				assert.Equal(t, []string{"player1", "player2"}, entities[0].PlayerIDs)
			},
		},
		{
			name: "find user by multiple ids successfully",
			query: duser.DomainQuery{
				Ids: []string{"user123", "user456"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1"},
				})
				if err != nil {
					t.Fatal(err)
				}
				err = repository.Save(ctx, duser.Entity{
					ID:        "user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
				// Save a user that should not appear
				err = repository.Save(ctx, duser.Entity{
					ID:        "user789",
					FirstName: "Bob",
					LastName:  "Johnson",
					PlayerIDs: []string{"player3"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				ids := []string{entities[0].ID, entities[1].ID}
				assert.Contains(t, ids, "user123")
				assert.Contains(t, ids, "user456")
				assert.NotContains(t, ids, "user789")
			},
		},
		{
			name: "find user by id and PlayerIDs successfully",
			query: duser.DomainQuery{
				Ids:       []string{"user123"},
				PlayerIDs: []string{"player1"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				// Save user that matches the query
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
				// Save another user with different ID to verify filtering works
				err = repository.Save(ctx, duser.Entity{
					ID:        "user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player3", "player4"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
				assert.Contains(t, entities[0].PlayerIDs, "player1")
			},
		},
		{
			name: "find user by id and PlayerIDs when PlayerIDs do not match",
			query: duser.DomainQuery{
				Ids:       []string{"user123"},
				PlayerIDs: []string{"player999"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				// Should be filtered out because PlayerIDs don't match
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "find user by multiple ids and PlayerIDs successfully",
			query: duser.DomainQuery{
				Ids:       []string{"user123", "user456"},
				PlayerIDs: []string{"player1"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
				err = repository.Save(ctx, duser.Entity{
					ID:        "user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player3", "player4"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				// Only user123 should be returned because it has player1
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
				assert.Contains(t, entities[0].PlayerIDs, "player1")
			},
		},
		{
			name: "find user by multiple ids and multiple PlayerIDs successfully",
			query: duser.DomainQuery{
				Ids:       []string{"user123", "user456"},
				PlayerIDs: []string{"player1", "player3"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				})
				if err != nil {
					t.Fatal(err)
				}
				err = repository.Save(ctx, duser.Entity{
					ID:        "user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player3", "player4"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				// Both users should be returned because user123 has player1 and user456 has player3
				assert.Len(t, entities, 2)
				ids := []string{entities[0].ID, entities[1].ID}
				assert.Contains(t, ids, "user123")
				assert.Contains(t, ids, "user456")
			},
		},
		{
			name: "find user that does not exist",
			query: duser.DomainQuery{
				Ids: []string{"nonexistent"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				// No users saved
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name:  "find with no criteria returns empty",
			query: duser.DomainQuery{},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "find with only PlayerIDs returns empty",
			query: duser.DomainQuery{
				PlayerIDs: []string{"player1"},
			},
			on: func(t *testing.T, repository duser.Repository) {
				err := repository.Save(ctx, duser.Entity{
					ID:        "user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1"},
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 0)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			testCase.on(t, repository)

			// when
			entities, err := repository.Find(ctx, testCase.query)

			// then
			testCase.assertions(t, entities, err)
		})
	}
}
