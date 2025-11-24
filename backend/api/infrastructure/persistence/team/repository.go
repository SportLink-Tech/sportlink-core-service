package team

import (
	"context"
	"fmt"
	"sportlink/api/domain/team"

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
		keyCond = expression.KeyAnd(keyCond, expression.KeyBeginsWith(expression.Key("Id"), query.Name))
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return []team.Entity{}, err
	}

	resp, err := repo.dbClient.Query(context.TODO(), &dynamodb.QueryInput{
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
		results = append(results, dto.ToDomain())
	}

	// Return empty slice if no results found
	if results == nil {
		return []team.Entity{}, nil
	}

	return results, nil
}

func From(entity team.Entity) (Dto, error) {
	if entity.Name == "" {
		return Dto{}, fmt.Errorf("ID could not be empty")
	}

	return Dto{
		EntityId: "Entity#Team",
		Id:       entity.Name,
		Category: int(entity.Category),
		Sport:    string(entity.Sport),
	}, nil
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

	// Combine all filters with AND
	if len(filters) > 0 {
		combinedFilter := filters[0]
		for i := 1; i < len(filters); i++ {
			combinedFilter = expression.And(combinedFilter, filters[i])
		}
		*builder = builder.WithFilter(combinedFilter)
	}
}
