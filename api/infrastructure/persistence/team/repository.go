package team

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"sportlink/api/domain/team"
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

// TODO implement me
func (repo *DynamoDBRepository) Save(entity team.Entity) error {
	return nil
}

// TODO implement me
func (repo *DynamoDBRepository) Find(query team.DomainQuery) ([]team.Entity, error) {
	return make([]team.Entity, 0), nil
}
