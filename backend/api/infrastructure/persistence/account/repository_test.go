package account_test

import (
	"context"
	"errors"
	daccount "sportlink/api/domain/account"
	"sportlink/api/infrastructure/persistence/account"
	amocks "sportlink/mocks/api/infrastructure/persistence/account"
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
		entity     daccount.Entity
		setupMock  func(*amocks.DynamoDBClientInterface, daccount.Entity)
		assertions func(*testing.T, error)
	}{
		{
			name: "given valid account when saving then saves successfully",
			entity: daccount.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity daccount.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
					if input.TableName == nil || *input.TableName != "SportLinkCore" {
						return false
					}
					// Verify it contains the email
					var savedDto account.Dto
					_ = attributevalue.UnmarshalMap(input.Item, &savedDto)
					return savedDto.Email == entity.Email
				})).Return(&dynamodb.PutItemOutput{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "given account with empty email when saving then returns error",
			entity: daccount.Entity{
				Email:    "",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity daccount.Entity) {
				// DTO creation should fail, so PutItem won't be called
			},
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "email could not be empty")
			},
		},
		{
			name: "given account when PutItem fails then returns error",
			entity: daccount.Entity{
				Email:    "test@example.com",
				Nickname: "testuser",
				Password: "ValidP@ssw0rd123",
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity daccount.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockClient := amocks.NewDynamoDBClientInterface(t)
			tc.setupMock(mockClient, tc.entity)

			repository := account.NewRepositoryWithInterface(mockClient, "SportLinkCore")

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
		query      daccount.DomainQuery
		setupMock  func(*amocks.DynamoDBClientInterface, daccount.DomainQuery)
		assertions func(*testing.T, []daccount.Entity, error)
	}{
		{
			name: "given single email when finding then returns account",
			query: daccount.DomainQuery{
				Emails: []string{"test@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				dto := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#test@example.com",
					Email:    "test@example.com",
					Nickname: "testuser",
					Password: "$2a$10$hashedpassword",
				}
				av, _ := attributevalue.MarshalMap(dto)
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.TableName != nil && *input.TableName == "SportLinkCore"
				})).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name: "given multiple emails when finding then returns multiple accounts",
			query: daccount.DomainQuery{
				Emails: []string{"user1@example.com", "user2@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				dto1 := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#user1@example.com",
					Email:    "user1@example.com",
					Nickname: "user1",
					Password: "$2a$10$hashed1",
				}
				av1, _ := attributevalue.MarshalMap(dto1)

				dto2 := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#user2@example.com",
					Email:    "user2@example.com",
					Nickname: "user2",
					Password: "$2a$10$hashed2",
				}
				av2, _ := attributevalue.MarshalMap(dto2)

				// First call for user1@example.com
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.TableName != nil && *input.TableName == "SportLinkCore"
				})).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av1},
				}, nil).Once()

				// Second call for user2@example.com
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.TableName != nil && *input.TableName == "SportLinkCore"
				})).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av2},
				}, nil).Once()
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 2)
				emails := []string{entities[0].Email, entities[1].Email}
				assert.Contains(t, emails, "user1@example.com")
				assert.Contains(t, emails, "user2@example.com")
			},
		},
		{
			name: "given single ID when finding then returns account",
			query: daccount.DomainQuery{
				Ids: []string{"EMAIL#test@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				dto := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#test@example.com",
					Email:    "test@example.com",
					Nickname: "testuser",
					Password: "$2a$10$hashedpassword",
				}
				av, _ := attributevalue.MarshalMap(dto)
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "EMAIL#test@example.com", entities[0].ID)
			},
		},
		{
			name: "given multiple IDs when finding then returns multiple accounts",
			query: daccount.DomainQuery{
				Ids: []string{"EMAIL#user1@example.com", "EMAIL#user2@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				dto1 := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#user1@example.com",
					Email:    "user1@example.com",
					Nickname: "user1",
					Password: "$2a$10$hashed1",
				}
				av1, _ := attributevalue.MarshalMap(dto1)

				dto2 := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#user2@example.com",
					Email:    "user2@example.com",
					Nickname: "user2",
					Password: "$2a$10$hashed2",
				}
				av2, _ := attributevalue.MarshalMap(dto2)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av1},
				}, nil).Once()

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av2},
				}, nil).Once()
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 2)
			},
		},
		{
			name: "given email and nickname when finding then returns filtered account",
			query: daccount.DomainQuery{
				Emails:    []string{"test@example.com"},
				Nicknames: []string{"testuser"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				// DynamoDB returns the account with matching email
				// The nickname filter is applied in-memory after the query
				dto := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#test@example.com",
					Email:    "test@example.com",
					Nickname: "testuser",
					Password: "$2a$10$hashed1",
				}
				av, _ := attributevalue.MarshalMap(dto)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name: "given email and nickname when nickname does not match then returns empty",
			query: daccount.DomainQuery{
				Emails:    []string{"test@example.com"},
				Nicknames: []string{"testuser"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				// DynamoDB returns an account with matching email but different nickname
				// The nickname filter should filter it out
				dto := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#test@example.com",
					Email:    "test@example.com",
					Nickname: "differentuser",
					Password: "$2a$10$hashed1",
				}
				av, _ := attributevalue.MarshalMap(dto)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				// Should be filtered out because nickname doesn't match
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given email and multiple nicknames when finding then returns filtered account",
			query: daccount.DomainQuery{
				Emails:    []string{"test@example.com"},
				Nicknames: []string{"testuser", "anotheruser"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				// DynamoDB can only return one item per email (because ID is unique)
				// This test verifies that when the nickname matches one of the filters, it's returned
				dto := account.Dto{
					EntityId: "Entity#Account",
					Id:       "EMAIL#test@example.com",
					Email:    "test@example.com",
					Nickname: "testuser",
					Password: "$2a$10$hashed1",
				}
				av, _ := attributevalue.MarshalMap(dto)

				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{av},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 1)
				assert.Equal(t, "test@example.com", entities[0].Email)
				assert.Equal(t, "testuser", entities[0].Nickname)
			},
		},
		{
			name:  "given no criteria when finding then returns empty",
			query: daccount.DomainQuery{},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				// No DynamoDB calls should be made
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given both Emails and Ids when finding then returns error",
			query: daccount.DomainQuery{
				Emails: []string{"test@example.com"},
				Ids:    []string{"EMAIL#test@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				// No DynamoDB calls should be made
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot use both Emails and Ids in query")
				assert.Len(t, entities, 0)
			},
		},
		{
			name: "given Query fails when finding then returns error",
			query: daccount.DomainQuery{
				Emails: []string{"test@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				mockClient.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "given Query returns empty when finding then returns empty slice",
			query: daccount.DomainQuery{
				Emails: []string{"nonexistent@example.com"},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query daccount.DomainQuery) {
				mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{},
				}, nil)
			},
			assertions: func(t *testing.T, entities []daccount.Entity, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockClient := amocks.NewDynamoDBClientInterface(t)
			tc.setupMock(mockClient, tc.query)

			repository := account.NewRepositoryWithInterface(mockClient, "SportLinkCore")

			// when
			entities, err := repository.Find(context.Background(), tc.query)

			// then
			tc.assertions(t, entities, err)
			mockClient.AssertExpectations(t)
		})
	}
}
