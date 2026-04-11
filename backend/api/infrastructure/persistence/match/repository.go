package match

import (
	"context"
	"fmt"
	"sportlink/api/domain/match"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClientInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
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

// Save writes two DynamoDB items — one per participant account — so both can
// efficiently list their matches by querying their own partition key.
func (repo *RepositoryAdapter) Save(ctx context.Context, entity match.Entity) error {
	localDto, visitorDto := fromEntity(entity)

	for _, dto := range []Dto{localDto, visitorDto} {
		av, err := attributevalue.MarshalMap(dto)
		if err != nil {
			return fmt.Errorf("failed to marshal match dto: %w", err)
		}
		_, err = repo.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(repo.tableName),
			Item:      av,
		})
		if err != nil {
			return fmt.Errorf("failed to save match record: %w", err)
		}
	}
	return nil
}

// Find returns all matches for the given account, optionally filtered by status.
func (repo *RepositoryAdapter) Find(ctx context.Context, query match.DomainQuery) ([]match.Entity, error) {
	keyCond := expression.KeyEqual(
		expression.Key("EntityId"),
		expression.Value(entityIDPrefix(query.AccountID)),
	)

	builder := expression.NewBuilder().WithKeyCondition(keyCond)

	if len(query.Statuses) > 0 {
		var vals []expression.OperandBuilder
		for _, s := range query.Statuses {
			vals = append(vals, expression.Value(s.String()))
		}
		builder = builder.WithFilter(expression.Name("Status").In(vals[0], vals[1:]...))
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	var entities []match.Entity
	var lastKey map[string]types.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			TableName:                 aws.String(repo.tableName),
			KeyConditionExpression:    expr.KeyCondition(),
			FilterExpression:          expr.Filter(),
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

		items, err := unmarshalItems(resp.Items)
		if err != nil {
			return nil, err
		}
		entities = append(entities, items...)

		if resp.LastEvaluatedKey == nil {
			break
		}
		lastKey = resp.LastEvaluatedKey
	}

	if entities == nil {
		return []match.Entity{}, nil
	}
	return entities, nil
}

// FindByID returns a single match by ID, scoped to one of its participant accounts.
func (repo *RepositoryAdapter) FindByID(ctx context.Context, accountID, matchID string) (*match.Entity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"EntityId": entityIDPrefix(accountID),
		"Id":       matchID,
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

	var dto Dto
	if err := attributevalue.UnmarshalMap(resp.Item, &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal match: %w", err)
	}
	entity := dto.ToDomain()
	return &entity, nil
}

func unmarshalItems(items []map[string]types.AttributeValue) ([]match.Entity, error) {
	entities := make([]match.Entity, 0, len(items))
	for _, item := range items {
		var dto Dto
		if err := attributevalue.UnmarshalMap(item, &dto); err != nil {
			return nil, fmt.Errorf("failed to unmarshal match: %w", err)
		}
		entities = append(entities, dto.ToDomain())
	}
	return entities, nil
}
