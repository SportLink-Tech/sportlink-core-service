package dynamodb_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sportlink/api/infrastructure/persistence/dynamodb"
	"sportlink/api/infrastructure/persistence/player"
	"sportlink/dev/testcontainer"
	"testing"
	"time"
)

func TestDynamoDb_Save(t *testing.T) {

	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := dynamodb.NewDynamoDb[player.Dto](*dynamoDbClient, "SportLinkCore")

	testCases := []struct {
		name       string
		itemToSave player.Dto
		assertions func(t *testing.T, err error)
	}{
		{
			name: "save an item successfully",
			itemToSave: player.Dto{
				EntityId: "EntityId#ItemDto",
				Id:       "1234",
				Category: 1,
				Sport:    "football",
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "save an item without id must fail",
			itemToSave: player.Dto{
				EntityId: "EntityId#ItemDto",
				Category: 1,
				Sport:    "football",
			},
			assertions: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "api error ValidationException: One or more parameter values are not valid. The AttributeValue for a key attribute cannot contain an empty string value.")
				assert.NotNil(t, err)
			},
		},
		{
			name: "save an item successfully with pk and sk",
			itemToSave: player.Dto{
				EntityId: "EntityId#ItemDto",
				Id:       "1234",
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := repository.Save(ctx, testCase.itemToSave)

			// then
			testCase.assertions(t, err)
		})
	}
}

func TestDynamoDb_SaveAll(t *testing.T) {

	ctx := context.Background()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	dynamoDbClient := testcontainer.GetDynamoDbClient(t, container, ctx)
	repository := dynamodb.NewDynamoDb[player.Dto](*dynamoDbClient, "SportLinkCore")
	time.Sleep(5 * time.Second)

	testCases := []struct {
		name        string
		itemsToSave []player.Dto
		batchSize   int
		assertions  func(t *testing.T, output dynamodb.SaveAllOutput, err error)
	}{
		{
			name: "save all items successfully with an small batch size successfully",
			itemsToSave: []player.Dto{
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1234",
					Category: 1,
					Sport:    "football",
				},
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1235",
					Category: 1,
					Sport:    "football",
				},
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1236",
					Category: 1,
					Sport:    "football",
				},
			},
			batchSize: 2,
			assertions: func(t *testing.T, output dynamodb.SaveAllOutput, err error) {
				assert.Equal(t, output.ProcessedItems, 3)
				assert.Equal(t, output.UnprocessedItems, 0)
				assert.Nil(t, err)
			},
		},
		{
			name: "save all items successfully with a batch size successfully",
			itemsToSave: []player.Dto{
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1234",
					Category: 1,
					Sport:    "football",
				},
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1235",
					Category: 1,
					Sport:    "football",
				},
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1236",
					Category: 1,
					Sport:    "football",
				},
			},
			batchSize: 25,
			assertions: func(t *testing.T, output dynamodb.SaveAllOutput, err error) {
				assert.Equal(t, output.ProcessedItems, 3)
				assert.Equal(t, output.UnprocessedItems, 0)
				assert.Nil(t, err)
			},
		},
		{
			name:        "try to save items with an empty list must not fail",
			itemsToSave: []player.Dto{},
			batchSize:   25,
			assertions: func(t *testing.T, output dynamodb.SaveAllOutput, err error) {
				assert.Equal(t, output.ProcessedItems, 0)
				assert.Equal(t, output.UnprocessedItems, 0)
				assert.Nil(t, err)
			},
		},
		{
			name: "when the final batch could not be saved then the information about the previous one must be retrieved",
			itemsToSave: []player.Dto{
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1234",
					Category: 1,
					Sport:    "football",
				},
				{
					EntityId: "EntityId#ItemDto",
					Id:       "1235",
					Category: 1,
					Sport:    "football",
				},
				// this item will be contained in a different batch
				{
					Id:       "1236",
					Category: 1,
					Sport:    "football",
				},
			},
			batchSize: 2,
			assertions: func(t *testing.T, output dynamodb.SaveAllOutput, err error) {
				assert.Contains(t, err.Error(), "api error ValidationException: One or more parameter values are not valid. The AttributeValue for a key attribute cannot contain an empty string value. Key: EntityId")
				assert.Equal(t, output.ProcessedItems, 2)
				assert.NotNil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			output, err := repository.SaveAll(ctx, testCase.itemsToSave, testCase.batchSize)

			// then
			testCase.assertions(t, output, err)
		})
	}
}
