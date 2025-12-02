package matchannouncement

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchannouncement"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type RepositoryAdapter struct {
	dbClient  *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) matchannouncement.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(entity matchannouncement.Entity) error {
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

func (repo *RepositoryAdapter) Find(query matchannouncement.DomainQuery) ([]matchannouncement.Entity, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#MatchAnnouncement"))

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return []matchannouncement.Entity{}, err
	}

	resp, err := repo.dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	if err != nil {
		return []matchannouncement.Entity{}, err
	}

	var results []matchannouncement.Entity
	for _, item := range resp.Items {
		var dto Dto
		err = attributevalue.UnmarshalMap(item, &dto)
		if err != nil {
			return []matchannouncement.Entity{}, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		results = append(results, dto.ToDomain())
	}

	// Return empty slice if no results found
	if results == nil {
		return []matchannouncement.Entity{}, nil
	}

	return results, nil
}

func From(entity matchannouncement.Entity) (Dto, error) {
	// Get the timezone from location
	tz := entity.Location.GetTimezone()

	// Convert times to Unix timestamps using the location's timezone
	day := entity.Day.In(tz).Unix()
	startTime := entity.TimeSlot.StartTime.In(tz).Unix()
	endTime := entity.TimeSlot.EndTime.In(tz).Unix()
	createdAt := entity.CreatedAt.In(tz).Unix()

	// Calculate ExpiresAt: 30 days from creation date
	expiresAt := entity.CreatedAt.In(tz).AddDate(0, 0, 30).Unix()

	// Extract category range data
	var categories []int
	var minLevel, maxLevel int

	switch entity.AdmittedCategories.Type {
	case matchannouncement.RangeTypeSpecific:
		categories = make([]int, len(entity.AdmittedCategories.Categories))
		for i, c := range entity.AdmittedCategories.Categories {
			categories[i] = int(c)
		}
	case matchannouncement.RangeTypeGreaterThan:
		minLevel = int(entity.AdmittedCategories.MinLevel)
	case matchannouncement.RangeTypeLessThan:
		maxLevel = int(entity.AdmittedCategories.MaxLevel)
	case matchannouncement.RangeTypeBetween:
		minLevel = int(entity.AdmittedCategories.MinLevel)
		maxLevel = int(entity.AdmittedCategories.MaxLevel)
	}

	return Dto{
		EntityId:   "Entity#MatchAnnouncement",
		Id:         entity.ID,
		TeamName:   entity.TeamName,
		Sport:      string(entity.Sport),
		Day:        day,
		StartTime:  startTime,
		EndTime:    endTime,
		Country:    entity.Location.Country,
		Province:   entity.Location.Province,
		Locality:   entity.Location.Locality,
		RangeType:  string(entity.AdmittedCategories.Type),
		Categories: categories,
		MinLevel:   minLevel,
		MaxLevel:   maxLevel,
		Status:     entity.Status.String(),
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
	}, nil
}

func includeFilters(query matchannouncement.DomainQuery, builder *expression.Builder) {
	var filters []expression.ConditionBuilder

	// Filter by sports
	if len(query.Sports) > 0 {
		var sportValues []expression.OperandBuilder
		for _, s := range query.Sports {
			sportValues = append(sportValues, expression.Value(string(s)))
		}
		filters = append(filters, expression.Name("Sport").In(sportValues[0], sportValues[1:]...))
	}

	// Note: Category filtering is complex due to CategoryRange types (SPECIFIC, GREATER_THAN, etc.)
	// For now, we skip category filtering in DynamoDB and can apply it in-memory if needed
	// This is because DynamoDB doesn't easily support filtering on complex range logic

	// Filter by statuses
	if len(query.Statuses) > 0 {
		var statusValues []expression.OperandBuilder
		for _, s := range query.Statuses {
			statusValues = append(statusValues, expression.Value(s.String()))
		}
		filters = append(filters, expression.Name("Status").In(statusValues[0], statusValues[1:]...))
	}

	// Filter by date range
	if !query.FromDate.IsZero() {
		filters = append(filters, expression.Name("Day").GreaterThanEqual(expression.Value(query.FromDate.Unix())))
	}
	if !query.ToDate.IsZero() {
		filters = append(filters, expression.Name("Day").LessThanEqual(expression.Value(query.ToDate.Unix())))
	}

	// Filter by location
	if query.Location != nil {
		if query.Location.Country != "" {
			filters = append(filters, expression.Name("Country").Equal(expression.Value(query.Location.Country)))
		}
		if query.Location.Province != "" {
			filters = append(filters, expression.Name("Province").Equal(expression.Value(query.Location.Province)))
		}
		if query.Location.Locality != "" {
			filters = append(filters, expression.Name("Locality").Equal(expression.Value(query.Location.Locality)))
		}
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
