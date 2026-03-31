package matchrequest

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchrequest"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBClientInterface defines the interface for DynamoDB operations needed by the repository
type DynamoDBClientInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type RepositoryAdapter struct {
	dbClient  DynamoDBClientInterface
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) matchrequest.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func NewRepositoryWithInterface(client DynamoDBClientInterface, tableName string) matchrequest.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(ctx context.Context, entity matchrequest.Entity) error {
	dto := From(entity)

	av, err := attributevalue.MarshalMap(dto)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      av,
	})
	return err
}

const ownerAccountIDIndexName = "OwnerAccountId-index"

func (repo *RepositoryAdapter) Find(ctx context.Context, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	hasOwnerCriteria := len(query.OwnerAccountIDs) > 0
	hasIDCriteria := len(query.IDs) > 0

	switch {
	case hasOwnerCriteria && len(query.OwnerAccountIDs) > 1:
		return repo.findByMultipleOwners(ctx, query)
	case hasOwnerCriteria:
		return repo.findByOwnerAccountID(ctx, query.OwnerAccountIDs[0], query)
	case hasIDCriteria && len(query.IDs) > 1:
		return repo.findByMultipleIDs(ctx, query)
	case hasIDCriteria:
		return repo.findByPrimaryKey(ctx, query.IDs[0], query)
	default:
		return []matchrequest.Entity{}, nil
	}
}

func (repo *RepositoryAdapter) findByOwnerAccountID(ctx context.Context, ownerAccountID string, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("OwnerAccountId"), expression.Value(ownerAccountID))
	// Restrict to match requests only (the GSI is shared with other entity types)
	entityFilter := expression.Equal(expression.Name("EntityId"), expression.Value("Entity#MatchRequest"))

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	builder = includeFilters(query, builder).WithFilter(entityFilter)

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		IndexName:                 aws.String(ownerAccountIDIndexName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	return unmarshalItems(resp.Items)
}

func (repo *RepositoryAdapter) findByPrimaryKey(ctx context.Context, id string, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	keyCond := expression.KeyAnd(
		expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#MatchRequest")),
		expression.KeyEqual(expression.Key("Id"), expression.Value(id)),
	)

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	builder = includeFilters(query, builder)

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}

	return unmarshalItems(resp.Items)
}

func (repo *RepositoryAdapter) findByMultipleOwners(ctx context.Context, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	seen := make(map[string]bool)
	var all []matchrequest.Entity

	for _, ownerID := range query.OwnerAccountIDs {
		results, err := repo.findByOwnerAccountID(ctx, ownerID, query)
		if err != nil {
			return nil, err
		}
		for _, e := range results {
			if !seen[e.ID] {
				seen[e.ID] = true
				all = append(all, e)
			}
		}
	}

	if all == nil {
		return []matchrequest.Entity{}, nil
	}
	return all, nil
}

func (repo *RepositoryAdapter) findByMultipleIDs(ctx context.Context, query matchrequest.DomainQuery) ([]matchrequest.Entity, error) {
	seen := make(map[string]bool)
	var all []matchrequest.Entity

	for _, id := range query.IDs {
		results, err := repo.findByPrimaryKey(ctx, id, query)
		if err != nil {
			return nil, err
		}
		for _, e := range results {
			if !seen[e.ID] {
				seen[e.ID] = true
				all = append(all, e)
			}
		}
	}

	if all == nil {
		return []matchrequest.Entity{}, nil
	}
	return all, nil
}

func (repo *RepositoryAdapter) UpdateStatus(ctx context.Context, id string, ownerAccountID string, newStatus matchrequest.Status) error {
	update := expression.Set(expression.Name("Status"), expression.Value(newStatus.String()))
	cond := expression.And(
		expression.Equal(expression.Name("OwnerAccountId"), expression.Value(ownerAccountID)),
		expression.Equal(expression.Name("Status"), expression.Value(matchrequest.StatusPending.String())),
	)

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		WithCondition(cond).
		Build()
	if err != nil {
		return err
	}

	key := map[string]types.AttributeValue{
		"EntityId": &types.AttributeValueMemberS{Value: "Entity#MatchRequest"},
		"Id":       &types.AttributeValueMemberS{Value: id},
	}

	_, err = repo.dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(repo.tableName),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("failed to update match request status: %w", err)
	}

	return nil
}

func includeFilters(query matchrequest.DomainQuery, builder expression.Builder) expression.Builder {
	var filters []expression.ConditionBuilder

	if len(query.RequesterAccountIDs) > 0 {
		var values []expression.OperandBuilder
		for _, id := range query.RequesterAccountIDs {
			values = append(values, expression.Value(id))
		}
		filters = append(filters, expression.Name("RequesterAccountId").In(values[0], values[1:]...))
	}

	if len(query.Statuses) > 0 {
		var values []expression.OperandBuilder
		for _, s := range query.Statuses {
			values = append(values, expression.Value(s.String()))
		}
		filters = append(filters, expression.Name("Status").In(values[0], values[1:]...))
	}

	if len(filters) == 0 {
		return builder
	}

	combined := filters[0]
	for i := 1; i < len(filters); i++ {
		combined = expression.And(combined, filters[i])
	}
	return builder.WithFilter(combined)
}

func unmarshalItems(items []map[string]types.AttributeValue) ([]matchrequest.Entity, error) {
	entities := make([]matchrequest.Entity, 0, len(items))
	for _, item := range items {
		var dto Dto
		if err := attributevalue.UnmarshalMap(item, &dto); err != nil {
			return nil, fmt.Errorf("failed to unmarshal match request: %w", err)
		}
		entities = append(entities, dto.ToDomain())
	}
	return entities, nil
}
