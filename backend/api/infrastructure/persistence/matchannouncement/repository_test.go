package matchannouncement_test

import (
	"context"
	"errors"
	ddomain "sportlink/api/domain/matchannouncement"
	"sportlink/api/infrastructure/persistence/matchannouncement"
	amocks "sportlink/mocks/api/infrastructure/persistence/matchannouncement"
	"testing"
	"time"

	"sportlink/api/domain/common"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_Save(t *testing.T) {
	location := ddomain.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	timeSlot, _ := ddomain.NewTimeSlot(startTime, endTime)

	testCases := []struct {
		name       string
		entity     ddomain.Entity
		setupMock  func(*amocks.DynamoDBClientInterface, ddomain.Entity)
		assertions func(*testing.T, error)
	}{
		{
			name: "given valid match announcement when saving then saves successfully",
			entity: ddomain.NewMatchAnnouncement(
				"Thunder Strikers",
				common.Paddle,
				tomorrow,
				timeSlot,
				location,
				ddomain.NewSpecificCategories([]common.Category{5, 6, 7}),
				ddomain.StatusPending,
				time.Now().In(tz),
			),
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity ddomain.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
					if input.TableName == nil || *input.TableName != "SportLinkCore" {
						return false
					}
					// Verify it contains the team name
					var savedDto matchannouncement.Dto
					_ = attributevalue.UnmarshalMap(input.Item, &savedDto)
					return savedDto.TeamName == entity.TeamName
				})).Return(&dynamodb.PutItemOutput{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "given match announcement with GreaterThan category when saving then saves successfully",
			entity: ddomain.NewMatchAnnouncement(
				"Elite Team",
				common.Tennis,
				tomorrow,
				timeSlot,
				location,
				ddomain.NewGreaterThanCategory(5),
				ddomain.StatusPending,
				time.Now().In(tz),
			),
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity ddomain.Entity) {
				mockClient.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "given match announcement when PutItem fails then returns error",
			entity: ddomain.NewMatchAnnouncement(
				"Thunder Strikers",
				common.Paddle,
				tomorrow,
				timeSlot,
				location,
				ddomain.NewSpecificCategories([]common.Category{5, 6, 7}),
				ddomain.StatusPending,
				time.Now().In(tz),
			),
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, entity ddomain.Entity) {
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

			repository := matchannouncement.NewRepositoryWithInterface(mockClient, "SportLinkCore")

			// when
			err := repository.Save(context.Background(), tc.entity)

			// then
			tc.assertions(t, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	location := ddomain.NewLocation("Argentina", "Buenos Aires", "CABA")
	tz := location.GetTimezone()
	tomorrow := time.Now().In(tz).AddDate(0, 0, 1)
	startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tz)
	endTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tz)
	timeSlot, _ := ddomain.NewTimeSlot(startTime, endTime)

	testCases := []struct {
		name       string
		query      ddomain.DomainQuery
		setupMock  func(*amocks.DynamoDBClientInterface, ddomain.DomainQuery)
		assertions func(*testing.T, ddomain.Page, error)
	}{
		{
			name: "given query by sport when finding then returns matching announcements",
			query: ddomain.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity := ddomain.NewMatchAnnouncement(
					"Paddle Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto, _ := matchannouncement.From(entity)
				av, _ := attributevalue.MarshalMap(dto)

				// Mock countTotal query (no limit, returns all matching items)
				// countTotal makes queries in a loop until LastEvaluatedKey is nil
				// We need to match queries without limit OR with limit 0
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || (input.Limit != nil && *input.Limit == 0)
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query (with filters, uses batch fetching)
				// fetchWithFilters uses DynamoDBBatchSize (100) and may make multiple queries
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0 && *input.Limit == 100
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.Len(t, page.Entities, 1)
				assert.Equal(t, "Paddle Team A", page.Entities[0].TeamName)
				assert.Equal(t, common.Paddle, page.Entities[0].Sport)
				assert.Equal(t, 1, page.Total)
			},
		},
		{
			name: "given query by multiple statuses when finding then returns matching announcements",
			query: ddomain.DomainQuery{
				Statuses: []ddomain.Status{ddomain.StatusPending, ddomain.StatusConfirmed},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity1 := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto1, _ := matchannouncement.From(entity1)
				av1, _ := attributevalue.MarshalMap(dto1)

				entity2 := ddomain.NewMatchAnnouncement(
					"Team B",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusConfirmed,
					time.Now().In(tz),
				)
				dto2, _ := matchannouncement.From(entity2)
				av2, _ := attributevalue.MarshalMap(dto2)

				// Mock countTotal query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || *input.Limit == 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1, av2},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1, av2},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 2)
				assert.Equal(t, 2, page.Total)
			},
		},
		{
			name: "given query with limit and offset when finding then returns paginated results",
			query: ddomain.DomainQuery{
				Sports: []common.Sport{common.Paddle},
				Limit:  2,
				Offset: 1,
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity1 := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto1, _ := matchannouncement.From(entity1)
				av1, _ := attributevalue.MarshalMap(dto1)

				entity2 := ddomain.NewMatchAnnouncement(
					"Team B",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto2, _ := matchannouncement.From(entity2)
				av2, _ := attributevalue.MarshalMap(dto2)

				entity3 := ddomain.NewMatchAnnouncement(
					"Team C",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto3, _ := matchannouncement.From(entity3)
				av3, _ := attributevalue.MarshalMap(dto3)

				// Mock countTotal query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || *input.Limit == 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1, av2, av3},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query (should fetch enough items for offset)
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1, av2, av3},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.Len(t, page.Entities, 2)
				assert.Equal(t, 3, page.Total)
			},
		},
		{
			name: "given query with filters when finding then returns filtered results",
			query: ddomain.DomainQuery{
				Sports:   []common.Sport{common.Paddle},
				Statuses: []ddomain.Status{ddomain.StatusPending},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity := ddomain.NewMatchAnnouncement(
					"Paddle Pending Team",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto, _ := matchannouncement.From(entity)
				av, _ := attributevalue.MarshalMap(dto)

				// Mock countTotal query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || *input.Limit == 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query (with filters, uses batch fetching)
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 1)
				assert.Equal(t, common.Paddle, page.Entities[0].Sport)
				assert.Equal(t, ddomain.StatusPending, page.Entities[0].Status)
				assert.Equal(t, 1, page.Total)
			},
		},
		{
			name: "given query with date range when finding then returns filtered results",
			query: ddomain.DomainQuery{
				FromDate: tomorrow,
				ToDate:   tomorrow.AddDate(0, 0, 7),
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto, _ := matchannouncement.From(entity)
				av, _ := attributevalue.MarshalMap(dto)

				// Mock countTotal query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || *input.Limit == 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 1)
				assert.Equal(t, 1, page.Total)
			},
		},
		{
			name: "given query with location when finding then returns filtered results",
			query: ddomain.DomainQuery{
				Location: func() *ddomain.Location {
					loc := ddomain.NewLocation("Argentina", "Buenos Aires", "CABA")
					return &loc
				}(),
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto, _ := matchannouncement.From(entity)
				av, _ := attributevalue.MarshalMap(dto)

				// Mock countTotal query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit == nil || *input.Limit == 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)

				// Mock actual find query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 1)
				assert.Equal(t, 1, page.Total)
			},
		},
		{
			name:  "given empty query when finding then returns all results",
			query: ddomain.DomainQuery{},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto, _ := matchannouncement.From(entity)
				av, _ := attributevalue.MarshalMap(dto)

				// Mock both countTotal and fetchWithoutFilters queries (both without limit when query.Limit == 0)
				// When there are no filters and limit is 0, both queries have no limit
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					// Both countTotal and fetchWithoutFilters should query with no limit when query.Limit is 0
					// and with FilterExpression == nil (no filters)
					return (input.Limit == nil || *input.Limit == 0) && input.FilterExpression == nil
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 1)
				assert.Equal(t, 1, page.Total)
			},
		},
		{
			name: "given Query fails when finding then returns error",
			query: ddomain.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				mockClient.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
		{
			name: "given query when countTotal has pagination then handles multiple pages",
			query: ddomain.DomainQuery{
				Sports: []common.Sport{common.Paddle},
			},
			setupMock: func(mockClient *amocks.DynamoDBClientInterface, query ddomain.DomainQuery) {
				entity1 := ddomain.NewMatchAnnouncement(
					"Team A",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto1, _ := matchannouncement.From(entity1)
				av1, _ := attributevalue.MarshalMap(dto1)

				entity2 := ddomain.NewMatchAnnouncement(
					"Team B",
					common.Paddle,
					tomorrow,
					timeSlot,
					location,
					ddomain.NewSpecificCategories([]common.Category{5}),
					ddomain.StatusPending,
					time.Now().In(tz),
				)
				dto2, _ := matchannouncement.From(entity2)
				av2, _ := attributevalue.MarshalMap(dto2)

				lastKey := map[string]types.AttributeValue{
					"EntityId": &types.AttributeValueMemberS{Value: "Entity#MatchAnnouncement"},
					"Id":       &types.AttributeValueMemberS{Value: "some-id"},
				}

				// First countTotal query page
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return (input.Limit == nil || *input.Limit == 0) && input.ExclusiveStartKey == nil
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1},
					LastEvaluatedKey: lastKey,
				}, nil)

				// Second countTotal query page
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return (input.Limit == nil || *input.Limit == 0) && input.ExclusiveStartKey != nil
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av2},
					LastEvaluatedKey: nil,
				}, nil)

				// Actual find query
				mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
					return input.Limit != nil && *input.Limit > 0
				})).Return(&dynamodb.QueryOutput{
					Items:            []map[string]types.AttributeValue{av1, av2},
					LastEvaluatedKey: nil,
				}, nil)
			},
			assertions: func(t *testing.T, page ddomain.Page, err error) {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(page.Entities), 1)
				assert.Equal(t, 2, page.Total)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockClient := amocks.NewDynamoDBClientInterface(t)
			tc.setupMock(mockClient, tc.query)

			repository := matchannouncement.NewRepositoryWithInterface(mockClient, "SportLinkCore")

			// when
			page, err := repository.Find(context.Background(), tc.query)

			// then
			tc.assertions(t, page, err)
			mockClient.AssertExpectations(t)
		})
	}
}
