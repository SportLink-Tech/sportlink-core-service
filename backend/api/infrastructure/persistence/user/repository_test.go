package user_test

import (
	"context"
	"errors"
	duser "sportlink/api/domain/user"
	"sportlink/api/infrastructure/persistence/user"
	amocks "sportlink/mocks/api/infrastructure/persistence/user"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_Save(t *testing.T) {
	testCases := []struct {
		name       string
		entity     duser.Entity
		setupMock  func(*amocks.DynamoDBClientInterface, duser.Entity)
		assertions func(*testing.T, error)
	}{
		{
			name: "given valid user when saving then saves successfully",
			entity: duser.Entity{
				ID:        "user123",
				FirstName: "John",
				LastName:  "Doe",
				PlayerIDs: []string{"player1", "player2"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity duser.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
					if input.TableName == nil || *input.TableName != "SportLinkCore" {
						return false
					}
					// Verify it contains the ID
					var savedDto user.Dto
					_ = attributevalue.UnmarshalMap(input.Item, &savedDto)
					return savedDto.Id == "ID#"+entity.ID && savedDto.FirstName == entity.FirstName
				})).Return(&dynamodb.PutItemOutput{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "given user with empty ID when saving then returns error",
			entity: duser.Entity{
				ID:        "",
				FirstName: "John",
				LastName:  "Doe",
				PlayerIDs: []string{},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity duser.Entity) {
				// DTO creation should fail, so PutItem won't be called
			},
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "id could not be empty")
			},
		},
		{
			name: "given user when PutItem fails then returns error",
			entity: duser.Entity{
				ID:        "user123",
				FirstName: "John",
				LastName:  "Doe",
				PlayerIDs: []string{"player1"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity duser.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "given user with empty PlayerIDs when saving then saves successfully",
			entity: duser.Entity{
				ID:        "user123",
				FirstName: "John",
				LastName:  "Doe",
				PlayerIDs: []string{},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity duser.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockClient := amocks.NewDynamoDBClientInterface(t)
			tc.setupMock(mockClient, tc.entity)

			repository := user.NewRepositoryWithInterface(mockClient, "SportLinkCore")

			// when
			err := repository.Save(context.Background(), tc.entity)

			// then
			tc.assertions(t, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	testCases := []struct {
		name       string
		query      duser.DomainQuery
		setupMock  func(*amocks.DynamoDBClientInterface, duser.DomainQuery)
		assertions func(*testing.T, []duser.Entity, error)
	}{
		{
			name: "given single ID when finding then returns user",
			query: duser.DomainQuery{
				Ids: []string{"user123"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				dto := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				}
				av, _ := attributevalue.MarshalMap(dto)
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.TableName != nil && *input.TableName == "SportLinkCore"
				})).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
				assert.Equal(t, "John", entities[0].FirstName)
				assert.Equal(t, "Doe", entities[0].LastName)
				assert.Equal(t, []string{"player1", "player2"}, entities[0].PlayerIDs)
			},
		},
		{
			name: "given multiple IDs when finding then returns multiple users",
			query: duser.DomainQuery{
				Ids: []string{"user123", "user456"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				dto1 := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1"},
				}
				av1, _ := attributevalue.MarshalMap(dto1)

				dto2 := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player2"},
				}
				av2, _ := attributevalue.MarshalMap(dto2)

				// First call for user123 - use mock.Anything to match any query
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av1},
				}, nil).Once()

				// Second call for user456
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av2},
				}, nil).Once()
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 2)
				ids := []string{entities[0].ID, entities[1].ID}
				assert.Contains(t, ids, "user123")
				assert.Contains(t, ids, "user456")
			},
		},
		{
			name: "given ID and PlayerIDs when finding then returns filtered user",
			query: duser.DomainQuery{
				Ids:       []string{"user123"},
				PlayerIDs: []string{"player1"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				// DynamoDB returns the user with matching ID
				// The PlayerIDs filter is applied in-memory after the query
				dto := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				}
				av, _ := attributevalue.MarshalMap(dto)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
				// Should have player1 in PlayerIDs
				assert.Contains(t, entities[0].PlayerIDs, "player1")
			},
		},
		{
			name: "given ID and PlayerIDs when PlayerIDs do not match then returns empty",
			query: duser.DomainQuery{
				Ids:       []string{"user123"},
				PlayerIDs: []string{"player999"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				// DynamoDB returns a user with matching ID but different PlayerIDs
				// The PlayerIDs filter should filter it out
				dto := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				}
				av, _ := attributevalue.MarshalMap(dto)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				// Should be filtered out because PlayerIDs don't match
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given multiple IDs and PlayerIDs when finding then returns filtered users",
			query: duser.DomainQuery{
				Ids:       []string{"user123", "user456"},
				PlayerIDs: []string{"player1"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				dto1 := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1", "player2"},
				}
				av1, _ := attributevalue.MarshalMap(dto1)

				dto2 := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user456",
					FirstName: "Jane",
					LastName:  "Smith",
					PlayerIDs: []string{"player3", "player4"},
				}
				av2, _ := attributevalue.MarshalMap(dto2)

				// First call for user123
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av1},
				}, nil).Once()

				// Second call for user456
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av2},
				}, nil).Once()
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				// Only user123 should be returned because it has player1
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
			},
		},
		{
			name: "given empty IDs when finding then returns empty",
			query: duser.DomainQuery{
				Ids: []string{},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				// No Query should be called
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given Query fails when finding then returns error",
			query: duser.DomainQuery{
				Ids: []string{"user123"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				mockClient.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "given Query returns empty items when finding then returns empty",
			query: duser.DomainQuery{
				Ids: []string{"user123"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{},
				}, nil)
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given multiple IDs with duplicates when finding then returns deduplicated users",
			query: duser.DomainQuery{
				Ids: []string{"user123", "user123"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				dto := user.Dto{
					EntityId:  "Entity#User",
					Id:        "ID#user123",
					FirstName: "John",
					LastName:  "Doe",
					PlayerIDs: []string{"player1"},
				}
				av, _ := attributevalue.MarshalMap(dto)

				// Should be called twice (once for each ID in the query)
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil).Twice()
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.NoError(t, err)
				// Should deduplicate and return only one user
				assert.Len(t, entities, 1)
				assert.Equal(t, "user123", entities[0].ID)
			},
		},
		{
			name: "given unmarshal error when finding then returns error",
			query: duser.DomainQuery{
				Ids: []string{"user123"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query duser.DomainQuery) {
				// Return invalid item that will cause unmarshal error
				// Using BOOL type for Id field (should be string) - this will cause unmarshal to fail
				invalidItem := map[string]types.AttributeValue{
					"EntityId":  &types.AttributeValueMemberS{Value: "Entity#User"},
					"Id":        &types.AttributeValueMemberBOOL{Value: true}, // Should be string, not bool - will cause error
					"FirstName": &types.AttributeValueMemberS{Value: "John"},
					"LastName":  &types.AttributeValueMemberS{Value: "Doe"},
					"PlayerIDs": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				}
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{invalidItem},
				}, nil)
			},
			assertions: func(t *testing.T, entities []duser.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to unmarshal item")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockClient := amocks.NewDynamoDBClientInterface(t)
			tc.setupMock(mockClient, tc.query)

			repository := user.NewRepositoryWithInterface(mockClient, "SportLinkCore")

			// when
			entities, err := repository.Find(context.Background(), tc.query)

			// then
			tc.assertions(t, entities, err)
			mockClient.AssertExpectations(t)
		})
	}
}
