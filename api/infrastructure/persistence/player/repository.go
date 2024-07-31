package player

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"sportlink/api/domain/player"
)

type DynamoDBRepository struct {
	dbClient  *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *DynamoDBRepository) Save(entity player.Entity) error {
	dto, err := From(entity)
	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(dto)
	if err != nil {
		return err
	}

	_, err = repo.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})
	return err
}

func (repo *DynamoDBRepository) Find(query player.DomainQuery) ([]player.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#Player"))

	if query.ID != "" {
		keyCond = expression.KeyAnd(keyCond, expression.KeyEqual(expression.Key("Id"), expression.Value(query.ID)))
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	builder = includeFilters(query, builder)
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	var results []player.Entity
	resp, err := repo.dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Items {
		var entity player.Entity
		err = attributevalue.UnmarshalMap(item, &entity)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		results = append(results, entity)
	}

	return results, nil
}

func includeFilters(query player.DomainQuery, builder expression.Builder) expression.Builder {
	if query.Category != 0 || query.Sport != "" {
		var filter expression.ConditionBuilder
		filtersInitialized := false

		if query.Category != 0 {
			categoryFilter := expression.Name("Category").Equal(expression.Value(int(query.Category)))
			filter = categoryFilter
			filtersInitialized = true
		}

		if query.Sport != "" {
			sportFilter := expression.Name("Sport").Equal(expression.Value(query.Sport))
			if filtersInitialized {
				filter = expression.And(filter, sportFilter)
			} else {
				filter = sportFilter
			}
		}

		builder = builder.WithFilter(filter)
	}
	return builder
}
