package team_test

import (
	"context"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
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
			entity: dteam.NewTeam(
				"Boca",
				common.L1,
				*common.NewStats(0, 0, 0),
				common.Football,
				[]player.Entity{},
			),

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
				err := repository.Save(ctx, dteam.NewTeam(
					"First Test Team",
					common.L6,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
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
				repository.Save(ctx, dteam.NewTeam(
					"Second Test A",
					common.L2,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))

				repository.Save(ctx, dteam.NewTeam(
					"Second Test B",
					common.L2,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))

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
				Name: "NonExistentTeam12345",
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
				repository.Save(ctx, dteam.NewTeam(
					"Team Cat Sport A",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Team Cat Sport B",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
				// Guardar equipo de Football L2 (no debe aparecer)
				repository.Save(ctx, dteam.NewTeam(
					"Team Cat Sport C",
					common.L2,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
				// Guardar equipo de Paddle L1 (no debe aparecer)
				repository.Save(ctx, dteam.NewTeam(
					"Team Cat Sport D",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
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
				repository.Save(ctx, dteam.NewTeam(
					"Multi Cat Team A",
					common.L5,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Multi Cat Team B",
					common.L7,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Multi Cat Team C",
					common.L3,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.GreaterOrEqual(t, len(entities), 2, "should have at least 2 entities")
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
				// Verify that Multi Cat Team C (L3) is NOT in results
				assert.False(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Multi Cat Team C"},
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
				repository.Save(ctx, dteam.NewTeam(
					"Sport Only Team A",
					common.L3,
					*common.NewStats(0, 0, 0),
					common.Tennis,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Sport Only Team B",
					common.L4,
					*common.NewStats(0, 0, 0),
					common.Tennis,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Sport Only Team C",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
			},
			assertions: func(t *testing.T, entities []dteam.Entity, err error) {
				assert.Nil(t, err)
				assert.GreaterOrEqual(t, len(entities), 2, "should have at least 2 entities")
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
				// Verify that Sport Only Team C (Football) is NOT in results
				assert.False(t, slice.Contains[dteam.Entity](
					entities,
					dteam.Entity{Name: "Sport Only Team C"},
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
				repository.Save(ctx, dteam.NewTeam(
					"Pattern Match Team",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Pattern No Match Category",
					common.L2,
					*common.NewStats(0, 0, 0),
					common.Football,
					[]player.Entity{},
				))
				repository.Save(ctx, dteam.NewTeam(
					"Pattern No Match Sport",
					common.L1,
					*common.NewStats(0, 0, 0),
					common.Paddle,
					[]player.Entity{},
				))
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
			entities, err := repository.Find(ctx, testCase.query)

			// then
			testCase.assertions(t, entities, err)
		})
	}
}
