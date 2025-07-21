package team

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"sportlink/api/domain/team"
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

func (repo *RepositoryAdapter) Save(entity team.Entity) error {
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

func (repo *RepositoryAdapter) Find(query team.DomainQuery) ([]team.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#Team"))

	if query.Name != "" {
		keyCond = expression.KeyAnd(keyCond, expression.KeyEqual(expression.Key("Id"), expression.Value(query.Name)))
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	var results []team.Entity
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
		var dto Dto
		err = attributevalue.UnmarshalMap(item, &dto)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		results = append(results, dto.ToDomain())
	}

	return results, nil
}

func From(entity team.Entity) (Dto, error) {
	if entity.Name == "" {
		return Dto{}, fmt.Errorf("Id could not be empty")
	}

	return Dto{
		EntityId: "Entity#Team",
		Id:       entity.Name,
		Category: int(entity.Category),
		Sport:    string(entity.Sport),
	}, nil
}

func includeFilters(query team.DomainQuery, builder *expression.Builder) {
	if len(query.Categories) > 0 {
		var categoryValues []expression.OperandBuilder
		for _, c := range query.Categories {
			categoryValues = append(categoryValues, expression.Value(int(c)))
		}

		filter := expression.Name("Category").In(categoryValues[0], categoryValues[1:]...)
		*builder = builder.WithFilter(filter)
	}
}
