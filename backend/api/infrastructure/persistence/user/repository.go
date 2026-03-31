package user

import (
	"context"
	"fmt"
	"sportlink/api/domain/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClientInterface defines the interface for DynamoDB operations needed by the repository
type DynamoDBClientInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type RepositoryAdapter struct {
	dbClient  DynamoDBClientInterface
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) user.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

// NewRepositoryWithInterface allows injecting a mock client for testing
func NewRepositoryWithInterface(client DynamoDBClientInterface, tableName string) user.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(ctx context.Context, entity user.Entity) error {
	dto, err := From(entity)
	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(dto)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})
	return err
}

func (repo *RepositoryAdapter) Find(ctx context.Context, query user.DomainQuery) ([]user.Entity, error) {
	// We only support Query operations (no Scan)
	// Must have at least Ids to build key condition
	hasIdCriteria := len(query.Ids) > 0

	if !hasIdCriteria {
		// No valid key condition criteria, return empty
		return []user.Entity{}, nil
	}

	// For multiple IDs, we need to make multiple queries and combine results
	if len(query.Ids) > 1 {
		return repo.findWithMultipleIds(ctx, query)
	}

	// Single ID or empty (handled by findWithQuery)
	return repo.findWithQuery(ctx, query)
}

func (repo *RepositoryAdapter) findWithMultipleIds(ctx context.Context, query user.DomainQuery) ([]user.Entity, error) {
	var allResults []user.Entity
	seenIds := make(map[string]bool)

	// Query each ID separately and combine results
	for _, id := range query.Ids {
		// Create a query for this single ID
		singleIdQuery := query
		singleIdQuery.Ids = []string{id}

		results, err := repo.findWithQuery(ctx, singleIdQuery)
		if err != nil {
			return []user.Entity{}, err
		}

		// Add results, avoiding duplicates
		for _, entity := range results {
			if !seenIds[entity.ID] {
				seenIds[entity.ID] = true
				allResults = append(allResults, entity)
			}
		}
	}

	// Apply playerIDs filters if present
	if len(query.PlayerIDs) > 0 {
		var filteredResults []user.Entity
		for _, entity := range allResults {
			if matchesPlayerIDsFilter(entity, query) {
				filteredResults = append(filteredResults, entity)
			}
		}
		return filteredResults, nil
	}

	return allResults, nil
}

func (repo *RepositoryAdapter) findWithQuery(ctx context.Context, query user.DomainQuery) ([]user.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#User"))

	// Build key condition based on query
	// We expect exactly 0 or 1 ID at this point (multiple IDs are handled by findWithMultipleIds)
	if len(query.Ids) == 1 {
		// Search by single ID with prefix
		idValue := fmt.Sprintf("ID#%s", query.Ids[0])
		keyCond = expression.KeyAnd(keyCond, expression.KeyEqual(expression.Key("Id"), expression.Value(idValue)))
	} else {
		// No valid key condition, return empty
		return []user.Entity{}, nil
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return []user.Entity{}, err
	}

	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return []user.Entity{}, err
	}

	var results []user.Entity
	for _, item := range resp.Items {
		var dto Dto
		err = attributevalue.UnmarshalMap(item, &dto)
		if err != nil {
			return []user.Entity{}, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		entity := dto.ToDomain()

		// Apply in-memory filters for playerIDs if needed
		if !matchesPlayerIDsFilter(entity, query) {
			continue
		}

		results = append(results, entity)
	}

	// Return empty slice if no results found
	if results == nil {
		return []user.Entity{}, nil
	}

	return results, nil
}

func matchesPlayerIDsFilter(entity user.Entity, query user.DomainQuery) bool {
	if len(query.PlayerIDs) > 0 {
		// Check if entity has any of the queried playerIDs
		for _, queriedPlayerID := range query.PlayerIDs {
			for _, entityPlayerID := range entity.PlayerIDs {
				if entityPlayerID == queriedPlayerID {
					return true
				}
			}
		}
		return false
	}
	return true
}

func includeFilters(query user.DomainQuery, builder *expression.Builder) {
	// Note: PlayerIDs filtering is complex (array contains check)
	// We handle it in-memory after fetching results
	// Multiple IDs are handled by making multiple queries,
	// so we don't need to add any DynamoDB filters here
}
