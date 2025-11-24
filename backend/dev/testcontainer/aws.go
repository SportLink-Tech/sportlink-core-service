package testcontainer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

// DynamoDbCfg
func GetDynamoDbClient(t *testing.T, container testcontainers.Container, ctx context.Context) *dynamodb.Client {
	endpoint := GetContainerEndpoint(t, container, ctx)
	awsCfg := GetAwsConfig(t, ctx, endpoint)
	return dynamodb.NewFromConfig(awsCfg)
}

// Sqs
func GetSqsClient(t *testing.T, container testcontainers.Container, ctx context.Context) *sqs.Client {
	endpoint := GetContainerEndpoint(t, container, ctx)
	awsCfg := GetAwsConfig(t, ctx, endpoint)
	return sqs.NewFromConfig(awsCfg)
}

func ClearDynamoDbTable(t *testing.T, dynamoDbClient *dynamodb.Client, tableName string) {
	scanOutput, err := dynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("Failed to scan DynamoDB table: %s", err)
	}

	for _, item := range scanOutput.Items {
		_, err := dynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"EntityId": item["EntityId"],
				"Id":       item["Id"],
			},
		})
		if err != nil {
			t.Fatalf("Failed to delete item: %s", err)
		}
	}
}

func GetAwsConfig(t *testing.T, ctx context.Context, endpoint string) aws.Config {
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           fmt.Sprintf("%s", endpoint),
					SigningRegion: "us-east-1",
				}, nil
			})),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %s", err)
	}
	return awsConfig
}
