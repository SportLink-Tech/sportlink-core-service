package player_test

import (
	"context"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	iplayer "sportlink/api/infrastructure/persistence/player"
	"sportlink/dev/testcontainer"
	"testing"
)

const tableName = "SportLinkCore"

func TestDynamoDBRepository_Save(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := iplayer.NewDynamoDBRepository(dynamoDbClient, tableName)

	tests := []struct {
		name    string
		entity  player.Entity
		failure bool
	}{
		{
			name: "saving a new valid player",
			entity: player.Entity{
				ID:       "jorgejcabrera",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			failure: false,
		},
		{
			name: "saving a valid player without Category must not failed",
			entity: player.Entity{
				ID:    "jorgejcabrera",
				Sport: common.Paddle,
			},
			failure: false,
		},
		{
			name: "saving a player without id must failed",
			entity: player.Entity{
				ID:       "",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.Save(ctx, tt.entity)
			if (err != nil) != tt.failure {
				t.Errorf("it was an error: %v", err)
				return
			}
			testcontainer.ClearDynamoDbTable(t, dynamoDbClient, tableName)
		})
	}
}

func TestDynamoDBRepository_Find(t *testing.T) {
	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := iplayer.NewDynamoDBRepository(dynamoDbClient, tableName)

	tests := []struct {
		name                    string
		query                   player.DomainQuery
		savedPlayers            []player.Entity
		expectedAmountOfPlayers int
	}{
		{
			name: "finding a player by id",
			savedPlayers: []player.Entity{
				{
					ID:       "jorgejcabrera",
					Category: common.L1,
					Sport:    common.Paddle,
				},
			},
			query: player.DomainQuery{
				Id: "jorgejcabrera",
			},
			expectedAmountOfPlayers: 1,
		},
		{
			name: "missing player by id",
			savedPlayers: []player.Entity{
				{
					ID:       "jorge",
					Category: common.L1,
					Sport:    common.Paddle,
				},
			},
			query: player.DomainQuery{
				Id: "jorgejcabrera",
			},
			expectedAmountOfPlayers: 0,
		},
		{
			name: "finding all players by category",
			savedPlayers: []player.Entity{
				{
					ID:       "jorge",
					Category: common.L1,
					Sport:    common.Paddle,
				},
				{
					ID:       "javier",
					Category: common.L1,
					Sport:    common.Paddle,
				},
				{
					ID:       "cabrera",
					Category: common.L2,
					Sport:    common.Football,
				},
			},
			query: player.DomainQuery{
				Category: common.L1,
			},
			expectedAmountOfPlayers: 2,
		},
		{
			name: "finding all players by sport",
			savedPlayers: []player.Entity{
				{
					ID:       "ruben",
					Category: common.L1,
					Sport:    common.Tennis,
				},
				{
					ID:       "anastasio",
					Category: common.L1,
					Sport:    common.Tennis,
				},
				{
					ID:       "diaz",
					Category: common.L4,
					Sport:    common.Tennis,
				},
			},
			query: player.DomainQuery{
				Sport: common.Tennis,
			},
			expectedAmountOfPlayers: 3,
		},
		{
			name: "finding all players by sport and category",
			savedPlayers: []player.Entity{
				{
					ID:       "ruben",
					Category: common.L1,
					Sport:    common.Tennis,
				},
				{
					ID:       "anastasio",
					Category: common.L1,
					Sport:    common.Tennis,
				},
				{
					ID:       "diaz",
					Category: common.L4,
					Sport:    common.Tennis,
				},
			},
			query: player.DomainQuery{
				Sport:    common.Tennis,
				Category: common.L1,
			},
			expectedAmountOfPlayers: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, entity := range tt.savedPlayers {
				err := repository.Save(ctx, entity)
				if err != nil {
					t.Fatalf("failed to save entity: %v", err)
				}
			}
			players, err := repository.Find(ctx, tt.query)
			if err != nil {
				t.Fatalf("failed to find players: %v", err)
			}
			if len(players) != tt.expectedAmountOfPlayers {
				t.Fatalf("failed to find players")
			}
			testcontainer.ClearDynamoDbTable(t, dynamoDbClient, tableName)
		})
	}
}
