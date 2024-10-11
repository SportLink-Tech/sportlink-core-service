package messaging_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sportlink/api/infrastructure/persistence/dynamodb"
	"sportlink/api/infrastructure/persistence/player"
	"sportlink/dev/testcontainer"
	"testing"
)

func TestDynamoDb_Save(t *testing.T) {

	ctx := context.TODO()
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
