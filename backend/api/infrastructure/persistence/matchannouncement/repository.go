package matchannouncement

import (
	"context"
	"fmt"
	"sportlink/api/domain/matchannouncement"

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
}

type RepositoryAdapter struct {
	dbClient  DynamoDBClientInterface
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) matchannouncement.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

// NewRepositoryWithInterface allows injecting a mock client for testing
func NewRepositoryWithInterface(client DynamoDBClientInterface, tableName string) matchannouncement.Repository {
	return &RepositoryAdapter{
		dbClient:  client,
		tableName: tableName,
	}
}

func (repo *RepositoryAdapter) Save(ctx context.Context, entity matchannouncement.Entity) error {
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

const (
	// DynamoDBBatchSize is the batch size used when fetching filtered results
	// DynamoDB limit applies before FilterExpression, so we fetch in batches
	// to ensure we get enough filtered items
	DynamoDBBatchSize = 100
)

func (repo *RepositoryAdapter) Find(ctx context.Context, query matchannouncement.DomainQuery) (matchannouncement.Page, error) {
	keyCond := expression.KeyEqual(expression.Key("EntityId"), expression.Value("Entity#MatchAnnouncement"))

	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	includeFilters(query, &builder)
	expr, err := builder.Build()
	if err != nil {
		return matchannouncement.Page{}, err
	}

	totalCount, err := repo.countTotal(ctx, query, expr)
	if err != nil {
		return matchannouncement.Page{}, err
	}

	queryInput := repo.buildQueryInput(expr)
	hasFilters := expr.Filter() != nil

	var results []matchannouncement.Entity
	if hasFilters {
		results, err = repo.fetchWithFilters(ctx, queryInput, query.Limit, query.Offset)
	} else {
		results, err = repo.fetchWithoutFilters(ctx, queryInput, query.Limit, query.Offset)
	}
	if err != nil {
		return matchannouncement.Page{}, err
	}

	results = applyPagination(results, query.Limit, query.Offset)

	if results == nil {
		results = []matchannouncement.Entity{}
	}

	return matchannouncement.Page{
		Entities: results,
		Total:    totalCount,
	}, nil
}

func (repo *RepositoryAdapter) buildQueryInput(expr expression.Expression) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
}

func (repo *RepositoryAdapter) fetchWithFilters(ctx context.Context, queryInput *dynamodb.QueryInput, limit, offset int) ([]matchannouncement.Entity, error) {
	itemsNeeded := calculateItemsNeeded(limit, offset)
	queryInput.Limit = aws.Int32(DynamoDBBatchSize)

	var results []matchannouncement.Entity
	var lastEvaluatedKey map[string]types.AttributeValue

	// If limit is 0, fetch all items that pass the filter (no limit)
	hasLimit := itemsNeeded > 0

	for {
		if lastEvaluatedKey != nil {
			queryInput.ExclusiveStartKey = lastEvaluatedKey
		}

		pageResults, nextKey, err := repo.fetchQueryPage(ctx, queryInput)
		if err != nil {
			return nil, err
		}

		results = append(results, pageResults...)

		// Break if no more pages or if we have enough items (when limit is set)
		if nextKey == nil || (hasLimit && len(results) >= itemsNeeded) {
			break
		}
		lastEvaluatedKey = nextKey
	}

	return results, nil
}

func (repo *RepositoryAdapter) fetchWithoutFilters(ctx context.Context, queryInput *dynamodb.QueryInput, limit, offset int) ([]matchannouncement.Entity, error) {
	applyDynamoDBLimit(limit, offset, queryInput)

	results, _, err := repo.fetchQueryPage(ctx, queryInput)
	return results, err
}

func (repo *RepositoryAdapter) fetchQueryPage(ctx context.Context, queryInput *dynamodb.QueryInput) ([]matchannouncement.Entity, map[string]types.AttributeValue, error) {
	resp, err := repo.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, nil, err
	}

	var results []matchannouncement.Entity
	for _, item := range resp.Items {
		dto, err := repo.unmarshalItem(item)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, dto.ToDomain())
	}

	return results, resp.LastEvaluatedKey, nil
}

func (repo *RepositoryAdapter) unmarshalItem(item map[string]types.AttributeValue) (Dto, error) {
	var dto Dto
	err := attributevalue.UnmarshalMap(item, &dto)
	if err != nil {
		return Dto{}, fmt.Errorf("failed to unmarshal item: %w", err)
	}
	return dto, nil
}

func calculateItemsNeeded(limit, offset int) int {
	itemsNeeded := limit
	if offset > 0 {
		itemsNeeded = limit + offset
	}
	return itemsNeeded
}

// countTotal counts the total number of entities matching the query (ignoring pagination)
func (repo *RepositoryAdapter) countTotal(ctx context.Context, query matchannouncement.DomainQuery, expr expression.Expression) (int, error) {
	// Query without limit to count all matching items
	// Note: For large datasets, this could be expensive. Consider using a separate count index if needed.
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	totalCount := 0
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		if lastEvaluatedKey != nil {
			queryInput.ExclusiveStartKey = lastEvaluatedKey
		}

		resp, err := repo.dbClient.Query(ctx, queryInput)
		if err != nil {
			return 0, err
		}

		totalCount += len(resp.Items)

		if resp.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = resp.LastEvaluatedKey
	}

	return totalCount, nil
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

// applyDynamoDBLimit sets the DynamoDB query limit, accounting for offset and filters
// Note: DynamoDB limit applies before FilterExpression, so we need to fetch more items
// to ensure we have enough results after filtering. We use a multiplier to account for
// potential filtered items, or fetch all if limit is small.
func applyDynamoDBLimit(limit, offset int, queryInput *dynamodb.QueryInput) {
	if limit > 0 {
		limitToFetch := limit
		if offset > 0 {
			limitToFetch = limit + offset
		}
		// DynamoDB limit applies BEFORE FilterExpression, so we need to fetch more
		// to account for items that will be filtered out. Use a multiplier if we have filters.
		// For small limits, we can be more aggressive with the multiplier.
		hasFilters := queryInput.FilterExpression != nil
		if hasFilters {
			// Multiply by 3-5x to account for filtered items, but cap at reasonable value
			multiplier := 5
			adjustedLimit := limitToFetch * multiplier
			// Cap at 1000 to avoid fetching too much
			if adjustedLimit > 1000 {
				adjustedLimit = 1000
			}
			limitToFetch = adjustedLimit
		}
		queryInput.Limit = aws.Int32(int32(limitToFetch))
	}
}

// applyPagination applies limit and offset to a slice of results
// Returns the paginated slice
func applyPagination(results []matchannouncement.Entity, limit, offset int) []matchannouncement.Entity {
	// Apply offset in memory
	if offset > 0 && len(results) > offset {
		results = results[offset:]
	}

	// Apply limit in memory (in case we fetched more than needed)
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results
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
