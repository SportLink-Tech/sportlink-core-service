package team

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sportlink/api/domain/common"
	dteam "sportlink/api/domain/team"
	"sportlink/api/infrastructure/persistence/team"
	"sportlink/dev/testcontainer"
	"sportlink/dev/utils/slice"
	"testing"
)

func Test_Save(t *testing.T) {

	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := team.NewRepository(dynamoDbClient, "SportLinkCore")
	testCases := []struct {
		name       string
		entity     dteam.Entity
		assertions func(t *testing.T, err error)
	}{
		{
			name: "save an item successfully",
			entity: dteam.Entity{
				Name:     "Boca",
				Category: common.L1,
				Sport:    common.Football,
			},

			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := repository.Save(testCase.entity)

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
	repository := team.NewRepository(dynamoDbClient, "SportLinkCore")
	testCases := []struct {
		name       string
		query      dteam.DomainQuery
		on         func(t *testing.T, repository dteam.Repository)
		assertions func(t *testing.T, entities []dteam.Entity, err error)
	}{
		{
			name: "find an item successfully",
			query: dteam.DomainQuery{
				Name: "Boca",
			},
			on: func(t *testing.T, repository dteam.Repository) {
				err := repository.Save(dteam.Entity{
					Name:     "Boca",
					Category: common.L1,
					Sport:    common.Football,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Boca"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find more than one item successfully",
			query: dteam.DomainQuery{
				Categories: []common.Category{common.L1},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				repository.Save(dteam.Entity{
					Name:     "Boca",
					Category: common.L1,
					Sport:    common.Football,
				})

				repository.Save(dteam.Entity{
					Name:     "Instituto",
					Category: common.L1,
					Sport:    common.Football,
				})

			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Boca"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Instituto"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find a item which does not exist",
			query: dteam.DomainQuery{
				Name: "River",
			},
			on: func(t *testing.T, repository dteam.Repository) {
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
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
			entities, err := repository.Find(testCase.query)

			// then
			testCase.assertions(t, entities, err)
		})
	}
}
