package account

import (
	"context"
	"fmt"
	"sportlink/api/domain/account"

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

func NewRepository(client *dynamodb.Client, tableName string) account.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

// NewRepositoryWithInterface allows injecting a mock client for testing
func NewRepositoryWithInterface(client DynamoDBClientInterface, tableName string) account.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(ctx context.Context, entity account.Entity) error {
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

func (repo *RepositoryAdapter) Find(ctx context.Context, query account.DomainQuery) ([]account.Entity, error) {
	// We only support Query operations (no Scan)
	// Must have at least Emails or Ids to build key condition
	hasEmailCriteria := len(query.Emails) > 0
	hasIdCriteria := len(query.Ids) > 0

	if !hasEmailCriteria && !hasIdCriteria {
		// No valid key condition criteria, return empty
		return []account.Entity{}, nil
	}

	// Cannot use both Emails and Ids at the same time
	if hasEmailCriteria && hasIdCriteria {
		return []account.Entity{}, fmt.Errorf("cannot use both Emails and Ids in query")
	}

	// For multiple emails/IDs, we need to make multiple queries and combine results
	if len(query.Emails) > 1 {
		return repo.findWithMultipleEmails(ctx, query)
	}
	if len(query.Ids) > 1 {
		return repo.findWithMultipleIds(ctx, query)
	}

	// Single email, single ID, or empty (handled by findWithQuery)
	return repo.findWithQuery(ctx, query)
}

func (repo *RepositoryAdapter) findWithMultipleEmails(ctx context.Context, query account.DomainQuery) ([]account.Entity, error) {
	var allResults []account.Entity
	seenEmails := make(map[string]bool)

	// Query each email separately and combine results
	for _, email := range query.Emails {
		// Create a query for this single email
		singleEmailQuery := query
		singleEmailQuery.Emails = []string{email}

		results, err := repo.findWithQuery(ctx, singleEmailQuery)
		if err != nil {
			return []account.Entity{}, err
		}

		// Add results, avoiding duplicates
		for _, entity := range results {
			if !seenEmails[entity.Email] {
				seenEmails[entity.Email] = true
				allResults = append(allResults, entity)
			}
		}
	}

	// Apply nickname filters if present (already applied in findWithQuery, but check again for consistency)
	if len(query.Nicknames) > 0 {
		var filteredResults []account.Entity
		for _, entity := range allResults {
			if matchesNicknameFilter(entity, query) {
				filteredResults = append(filteredResults, entity)
			}
		}
		return filteredResults, nil
	}

	return allResults, nil
}

func (repo *RepositoryAdapter) findWithMultipleIds(ctx context.Context, query account.DomainQuery) ([]account.Entity, error) {
	var allResults []account.Entity
	seenIds := make(map[string]bool)

	// Query each ID separately and combine results
	for _, id := range query.Ids {
		// Create a query for this single ID
		singleIdQuery := query
		singleIdQuery.Ids = []string{id}

		results, err := repo.findWithQuery(ctx, singleIdQuery)
		if err != nil {
			return []account.Entity{}, err
		}

		// Add results, avoiding duplicates
		for _, entity := range results {
			if !seenIds[entity.ID] {
				seenIds[entity.ID] = true
				allResults = append(allResults, entity)
			}
		}
	}

	// Apply nickname filters if present
	if len(query.Nicknames) > 0 {
		var filteredResults []account.Entity
		for _, entity := range allResults {
			if matchesNicknameFilter(entity, query) {
				filteredResults = append(filteredResults, entity)
			}
		}
		return filteredResults, nil
	}

	return allResults, nil
}

func (repo *RepositoryAdapter) findWithQuery(ctx context.Context, query account.DomainQuery) ([]account.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#Account"))

	// Build key condition based on query
	// We expect exactly 0 or 1 email/ID at this point (multiple emails/IDs are handled by findWithMultiple*)
	if len(query.Emails) == 1 {
		// Search by single email
		idValue := fmt.Sprintf("EMAIL#%s", query.Emails[0])
		keyCond = expression.KeyAnd(keyCond, expression.KeyEqual(expression.Key("Id"), expression.Value(idValue)))
	} else if len(query.Ids) == 1 {
		// Search by single ID
		keyCond = expression.KeyAnd(keyCond, expression.KeyEqual(expression.Key("Id"), expression.Value(query.Ids[0])))
	} else {
		// No valid key condition, return empty
		return []account.Entity{}, nil
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return []account.Entity{}, err
	}

	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return []account.Entity{}, err
	}

	var results []account.Entity
	for _, item := range resp.Items {
		var dto Dto
		err = attributevalue.UnmarshalMap(item, &dto)
		if err != nil {
			return []account.Entity{}, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		entity := dto.ToDomain()

		// Apply in-memory filters for nicknames if needed
		if !matchesNicknameFilter(entity, query) {
			continue
		}

		results = append(results, entity)
	}

	// Return empty slice if no results found
	if results == nil {
		return []account.Entity{}, nil
	}

	return results, nil
}

func matchesNicknameFilter(entity account.Entity, query account.DomainQuery) bool {
	if len(query.Nicknames) > 0 {
		for _, nickname := range query.Nicknames {
			if entity.Nickname == nickname {
				return true
			}
		}
		return false
	}
	return true
}

func includeFilters(query account.DomainQuery, builder *expression.Builder) {
	var filters []expression.ConditionBuilder

	// Filter by nicknames (when used with other key conditions)
	if len(query.Nicknames) > 0 {
		var nicknameValues []expression.OperandBuilder
		for _, nickname := range query.Nicknames {
			nicknameValues = append(nicknameValues, expression.Value(nickname))
		}
		filters = append(filters, expression.Name("Nickname").In(nicknameValues[0], nicknameValues[1:]...))
	}

	// Note: Multiple emails/IDs are handled by making multiple queries,
	// so we don't need to filter them here

	// Combine all filters with AND
	if len(filters) > 0 {
		combinedFilter := filters[0]
		for i := 1; i < len(filters); i++ {
			combinedFilter = expression.And(combinedFilter, filters[i])
		}
		*builder = builder.WithFilter(combinedFilter)
	}
}
