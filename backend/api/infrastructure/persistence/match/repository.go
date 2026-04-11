package match

import (
	"context"
	"fmt"
	"sportlink/api/domain/match"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClientInterface interface {
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
}

type RepositoryAdapter struct {
	dbClient  DynamoDBClientInterface
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) match.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func NewRepositoryWithInterface(client DynamoDBClientInterface, tableName string) match.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

// Save persists a match atomically as three DynamoDB items:
//   - one canonical record (source of truth for all mutable data)
//   - two immutable pointer records (one per participant account, for efficient listing)
func (repo *RepositoryAdapter) Save(ctx context.Context, entity match.Entity) error {
	canonical, localPtr, visitorPtr := fromEntity(entity)

	canonicalAV, err := attributevalue.MarshalMap(canonical)
	if err != nil {
		return fmt.Errorf("failed to marshal canonical match: %w", err)
	}

	localAV, err := attributevalue.MarshalMap(localPtr)
	if err != nil {
		return fmt.Errorf("failed to marshal local match account pointer: %w", err)
	}

	visitorAV, err := attributevalue.MarshalMap(visitorPtr)
	if err != nil {
		return fmt.Errorf("failed to marshal visitor match account pointer: %w", err)
	}

	_, err = repo.dbClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{Put: &types.Put{TableName: aws.String(repo.tableName), Item: canonicalAV}},
			{Put: &types.Put{TableName: aws.String(repo.tableName), Item: localAV}},
			{Put: &types.Put{TableName: aws.String(repo.tableName), Item: visitorAV}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to save match transaction: %w", err)
	}

	return nil
}

// Find lists all matches for a given account by:
//  1. Querying the immutable pointer records (EntityId = "Entity#MatchAccount#<accountId>")
//  2. BatchGetItem on the canonical records to get up-to-date match data
//
// An optional status filter is applied after fetching canonical records.
func (repo *RepositoryAdapter) Find(ctx context.Context, query match.DomainQuery) ([]match.Entity, error) {
	matchIDs, err := repo.findMatchIDsByAccount(ctx, query.AccountID)
	if err != nil {
		return nil, err
	}
	if len(matchIDs) == 0 {
		return []match.Entity{}, nil
	}

	entities, err := repo.batchGetCanonical(ctx, matchIDs)
	if err != nil {
		return nil, err
	}

	if len(query.Statuses) == 0 {
		return entities, nil
	}

	allowed := make(map[match.Status]struct{}, len(query.Statuses))
	for _, s := range query.Statuses {
		allowed[s] = struct{}{}
	}
	filtered := make([]match.Entity, 0, len(entities))
	for _, e := range entities {
		if _, ok := allowed[e.Status]; ok {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// FindByID returns the canonical match record directly by its ID.
func (repo *RepositoryAdapter) FindByID(ctx context.Context, accountID, matchID string) (*match.Entity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"EntityId": canonicalEntityID,
		"Id":       canonicalIDKey(matchID),
	})
	if err != nil {
		return nil, err
	}

	resp, err := repo.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, nil
	}

	var dto MatchDto
	if err := attributevalue.UnmarshalMap(resp.Item, &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match: %w", err)
	}
	entity := dto.ToDomain()
	return &entity, nil
}

// findMatchIDsByAccount queries the pointer records to get all match IDs for an account.
func (repo *RepositoryAdapter) findMatchIDsByAccount(ctx context.Context, accountID string) ([]string, error) {
	keyCond := expression.KeyEqual(
		expression.Key("EntityId"),
		expression.Value(matchAccountEntityID(accountID)),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	var matchIDs []string
	var lastKey map[string]types.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			TableName:                 aws.String(repo.tableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		}
		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		resp, err := repo.dbClient.Query(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			var ptr MatchAccountDto
			if err := attributevalue.UnmarshalMap(item, &ptr); err != nil {
				return nil, fmt.Errorf("failed to unmarshal match account pointer: %w", err)
			}
			matchIDs = append(matchIDs, strings.TrimPrefix(ptr.Id, "Match#"))
		}

		if resp.LastEvaluatedKey == nil {
			break
		}
		lastKey = resp.LastEvaluatedKey
	}

	return matchIDs, nil
}

// batchGetCanonical fetches canonical MatchDto records for a list of match IDs.
func (repo *RepositoryAdapter) batchGetCanonical(ctx context.Context, matchIDs []string) ([]match.Entity, error) {
	keys := make([]map[string]types.AttributeValue, 0, len(matchIDs))
	for _, id := range matchIDs {
		key, err := attributevalue.MarshalMap(map[string]string{
			"EntityId": canonicalEntityID,
			"Id":       canonicalIDKey(id),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal batch key: %w", err)
		}
		keys = append(keys, key)
	}

	resp, err := repo.dbClient.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			repo.tableName: {Keys: keys},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to batch get canonical matches: %w", err)
	}

	items := resp.Responses[repo.tableName]
	entities := make([]match.Entity, 0, len(items))
	for _, item := range items {
		var dto MatchDto
		if err := attributevalue.UnmarshalMap(item, &dto); err != nil {
			return nil, fmt.Errorf("failed to unmarshal canonical match: %w", err)
		}
		entities = append(entities, dto.ToDomain())
	}
	return entities, nil
}
