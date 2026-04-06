package team

import (
	"context"
	"fmt"
	"sportlink/api/domain/team"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type RepositoryAdapter struct {
	dbClient  *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) team.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(ctx context.Context, entity team.Entity) error {
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

func (repo *RepositoryAdapter) Update(ctx context.Context, oldID string, entity team.Entity) error {
	if err := repo.Save(ctx, entity); err != nil {
		return err
	}

	if oldID != entity.ID {
		key, err := attributevalue.MarshalMap(map[string]interface{}{
			"EntityId": "Entity#Team",
			"Id":       oldID,
		})
		if err != nil {
			return err
		}
		_, err = repo.dbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(repo.tableName),
			Key:       key,
		})
		return err
	}

	return nil
}

func (repo *RepositoryAdapter) Find(ctx context.Context, query team.DomainQuery) ([]team.Entity, error) {
	if query.OwnerAccountID != "" {
		return repo.findByOwner(ctx, query.OwnerAccountID)
	}

	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#Team"))

	if query.Name != "" && len(query.Sports) > 0 {
		// Search by the full ID format: SPORT#<sport>#NAME#<name>
		// Use the first sport for the key condition, filter the rest if needed
		idPrefix := fmt.Sprintf("SPORT#%s#NAME#%s", query.Sports[0], query.Name)
		keyCond = expression.KeyAnd(keyCond, expression.KeyBeginsWith(expression.Key("Id"), idPrefix))
	} else if query.Name != "" {
		// If only name specified, search by name pattern (requires filter, not key condition)
		// We'll search for any ID that starts with SPORT# and filter by name
		keyCond = expression.KeyAnd(keyCond, expression.KeyBeginsWith(expression.Key("Id"), "SPORT#"))
	} else if len(query.Sports) > 0 {
		// If only sports specified, search by sport prefix
		sportPrefix := fmt.Sprintf("SPORT#%s", query.Sports[0])
		keyCond = expression.KeyAnd(keyCond, expression.KeyBeginsWith(expression.Key("Id"), sportPrefix))
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return []team.Entity{}, err
	}

	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return []team.Entity{}, err
	}

	var results []team.Entity
	for _, item := range resp.Items {
		var dto Dto
		err = attributevalue.UnmarshalMap(item, &dto)
		if err != nil {
			return []team.Entity{}, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		entity := dto.ToDomain()

		// If only name is specified without sports, filter by name in memory
		// (DynamoDB Contains filter may not work as expected for partial matches)
		if query.Name != "" && len(query.Sports) == 0 {
			if !strings.Contains(entity.Name, query.Name) {
				continue
			}
		}

		results = append(results, entity)
	}

	// Return empty slice if no results found
	if results == nil {
		return []team.Entity{}, nil
	}

	return results, nil
}

func From(entity team.Entity) (Dto, error) {
	if entity.ID == "" {
		return Dto{}, fmt.Errorf("ID could not be empty")
	}

	return Dto{
		EntityId:       "Entity#Team",
		Id:             entity.ID,
		Name:           entity.Name,
		Category:       int(entity.Category),
		Sport:          string(entity.Sport),
		OwnerAccountId: entity.OwnerAccountID,
	}, nil
}

func (repo *RepositoryAdapter) findByOwner(ctx context.Context, ownerAccountID string) ([]team.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("OwnerAccountId"), expression.Value(ownerAccountID))
	filter := expression.Name("EntityId").Equal(expression.Value("Entity#Team"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithFilter(filter).Build()
	if err != nil {
		return []team.Entity{}, err
	}

	indexName := "OwnerAccountId-index"
	resp, err := repo.dbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return []team.Entity{}, err
	}

	var results []team.Entity
	for _, item := range resp.Items {
		var dto Dto
		if err := attributevalue.UnmarshalMap(item, &dto); err != nil {
			return []team.Entity{}, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		results = append(results, dto.ToDomain())
	}

	if results == nil {
		return []team.Entity{}, nil
	}
	return results, nil
}

func includeFilters(query team.DomainQuery, builder *expression.Builder) {
	var filters []expression.ConditionBuilder

	if len(query.Categories) > 0 {
		var categoryValues []expression.OperandBuilder
		for _, c := range query.Categories {
			categoryValues = append(categoryValues, expression.Value(int(c)))
		}
		filters = append(filters, expression.Name("Category").In(categoryValues[0], categoryValues[1:]...))
	}

	if len(query.Sports) > 0 {
		var sportValues []expression.OperandBuilder
		for _, c := range query.Sports {
			sportValues = append(sportValues, expression.Value(string(c)))
		}
		filters = append(filters, expression.Name("Sport").In(sportValues[0], sportValues[1:]...))
	}

	// Note: Name filtering without sports is done in memory after query
	// to avoid DynamoDB Contains filter limitations

	// Combine all filters with AND
	if len(filters) > 0 {
		combinedFilter := filters[0]
		for i := 1; i < len(filters); i++ {
			combinedFilter = expression.And(combinedFilter, filters[i])
		}
		*builder = builder.WithFilter(combinedFilter)
	}
}
