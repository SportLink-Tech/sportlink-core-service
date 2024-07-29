package player

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func (repo *DynamoDBRepository) Insert(entity player.Entity) error {
	dto := From(entity)
	av, err := attributevalue.MarshalMap(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}

	_, err = repo.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})
	return err
}

func (repo *DynamoDBRepository) Update(entity player.Entity) error {
	av, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}

	_, err = repo.dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &repo.tableName,
		Item:      av,
	})
	return err
}

func (repo *DynamoDBRepository) Find(query player.DomainQuery) ([]player.Entity, error) {
	// Similar a Insert/Update, implementar la consulta aquí
	return nil, nil // Implementar según necesidad
}
