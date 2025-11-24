package team

import (
	"context"
	"sportlink/api/domain/common"
	dteam "sportlink/api/domain/team"
	"sportlink/api/infrastructure/persistence/team"
	"sportlink/dev/testcontainer"
	"sportlink/dev/utils/slice"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Name: "First",
			},
			on: func(t *testing.T, repository dteam.Repository) {
				err := repository.Save(dteam.Entity{
					Name:     "First Test Team",
					Category: common.L6,
					Sport:    common.Paddle,
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
					dteam.Entity{Name: "First Test Team"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find more than one item successfully",
			query: dteam.DomainQuery{
				Categories: []common.Category{common.L2},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				repository.Save(dteam.Entity{
					Name:     "Second Test A",
					Category: common.L2,
					Sport:    common.Football,
				})

				repository.Save(dteam.Entity{
					Name:     "Second Test B",
					Category: common.L2,
					Sport:    common.Paddle,
				})

			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Second Test A"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Second Test B"},
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
		{
			name: "find teams by category and sport",
			query: dteam.DomainQuery{
				Sports:     []common.Sport{common.Football},
				Categories: []common.Category{common.L1},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				// Guardar equipos de Football L1
				repository.Save(dteam.Entity{
					Name:     "Team Cat Sport A",
					Category: common.L1,
					Sport:    common.Football,
				})
				repository.Save(dteam.Entity{
					Name:     "Team Cat Sport B",
					Category: common.L1,
					Sport:    common.Football,
				})
				// Guardar equipo de Football L2 (no debe aparecer)
				repository.Save(dteam.Entity{
					Name:     "Team Cat Sport C",
					Category: common.L2,
					Sport:    common.Football,
				})
				// Guardar equipo de Paddle L1 (no debe aparecer)
				repository.Save(dteam.Entity{
					Name:     "Team Cat Sport D",
					Category: common.L1,
					Sport:    common.Paddle,
				})
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Team Cat Sport A"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Team Cat Sport B"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find teams by multiple categories in same sport",
			query: dteam.DomainQuery{
				Sports:     []common.Sport{common.Paddle},
				Categories: []common.Category{common.L5, common.L7},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				repository.Save(dteam.Entity{
					Name:     "Multi Cat Team A",
					Category: common.L5,
					Sport:    common.Paddle,
				})
				repository.Save(dteam.Entity{
					Name:     "Multi Cat Team B",
					Category: common.L7,
					Sport:    common.Paddle,
				})
				repository.Save(dteam.Entity{
					Name:     "Multi Cat Team C",
					Category: common.L3,
					Sport:    common.Paddle,
				})
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Multi Cat Team A"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Multi Cat Team B"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find teams by sport only (all categories)",
			query: dteam.DomainQuery{
				Sports: []common.Sport{common.Tennis},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				repository.Save(dteam.Entity{
					Name:     "Sport Only Team A",
					Category: common.L3,
					Sport:    common.Tennis,
				})
				repository.Save(dteam.Entity{
					Name:     "Sport Only Team B",
					Category: common.L4,
					Sport:    common.Tennis,
				})
				repository.Save(dteam.Entity{
					Name:     "Sport Only Team C",
					Category: common.L1,
					Sport:    common.Football,
				})
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 2)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Sport Only Team A"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Sport Only Team B"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
			},
		},
		{
			name: "find teams by name pattern, sport and category",
			query: dteam.DomainQuery{
				Name:       "Pattern",
				Sports:     []common.Sport{common.Football},
				Categories: []common.Category{common.L1},
			},
			on: func(t *testing.T, repository dteam.Repository) {
				repository.Save(dteam.Entity{
					Name:     "Pattern Match Team",
					Category: common.L1,
					Sport:    common.Football,
				})
				repository.Save(dteam.Entity{
					Name:     "Pattern No Match Category",
					Category: common.L2,
					Sport:    common.Football,
				})
				repository.Save(dteam.Entity{
					Name:     "Pattern No Match Sport",
					Category: common.L1,
					Sport:    common.Paddle,
				})
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.Len(t, entities, 1)
				assert.True(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Pattern Match Team"},
					func(a, b dteam.Entity) bool {
						return a.Name == b.Name
					}))
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
