package player

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"testing"
)

func TestDynamoDBRepository_Save(t *testing.T) {
	ctx := context.Background()
	localstackContainerRequest := createLocalstackContainerRequest()
	container := createContainer(t, ctx, localstackContainerRequest)
	defer container.Terminate(ctx)
	dynamoDbClient := getDynamoDbClient(t, container, ctx)
	tableName := "SportLinkCore"
	repository := NewDynamoDBRepository(dynamoDbClient, tableName)

	tests := []struct {
		name    string
		entity  player.Entity
		failure bool
	}{
		{
			name: "saving a new valid player",
			entity: player.Entity{
				ID:       "jorgejcabrera",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			failure: false,
		},
		{
			name: "saving a valid player without Category must not failed",
			entity: player.Entity{
				ID:    "jorgejcabrera",
				Sport: common.Paddle,
			},
			failure: false,
		},
		{
			name: "saving a player without id must failed",
			entity: player.Entity{
				ID:       "",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.Save(tt.entity)
			if (err != nil) != tt.failure {
				t.Errorf("it was an error: %v", err)
				return
			}
			clearDynamoTable(t, dynamoDbClient, tableName)
		})
	}
}

func TestDynamoDBRepository_Find(t *testing.T) {
	ctx := context.Background()
	localstackContainerRequest := createLocalstackContainerRequest()
	container := createContainer(t, ctx, localstackContainerRequest)
	defer container.Terminate(ctx)
	dynamoDbClient := getDynamoDbClient(t, container, ctx)
	tableName := "SportLinkCore"
	repository := NewDynamoDBRepository(dynamoDbClient, tableName)

	tests := []struct {
		name        string
		savedEntity player.Entity
		query       player.DomainQuery
		found       bool
	}{
		{
			name: "finding a player by id",
			savedEntity: player.Entity{
				ID:       "jorgejcabrera",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			query: player.DomainQuery{
				ID: "jorgejcabrera",
			},
			found: true,
		},
		{
			name: "missing player by id",
			savedEntity: player.Entity{
				ID:       "jorge",
				Category: common.L1,
				Sport:    common.Paddle,
			},
			query: player.DomainQuery{
				ID: "jorgejcabrera",
			},
			found: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.Save(tt.savedEntity)
			if err != nil {
				t.Fatalf("failed to save entity: %v", err)
			}
			players, err := repository.Find(tt.query)
			if err != nil {
				t.Fatalf("failed to find players: %v", err)
			}
			if (len(players) >= 1 && !tt.found) || (len(players) == 0 && tt.found) {
				t.Fatalf("failed to find players")
			}
			clearDynamoTable(t, dynamoDbClient, tableName)
		})
	}
}

// TODO move all this code to dev package
func getDynamoDbClient(t *testing.T, container testcontainers.Container, ctx context.Context) *dynamodb.Client {
	endpoint := getContainerEndpoint(t, container, ctx)
	awsCfg := getAwsConfig(t, ctx, endpoint)
	return dynamodb.NewFromConfig(awsCfg)
}

func getContainerEndpoint(t *testing.T, container testcontainers.Container, ctx context.Context) string {
	endpoint, err := container.PortEndpoint(ctx, "4566/tcp", "http")
	if err != nil {
		t.Fatalf("Failed to get container endpoint: %s", err)
	}
	return endpoint
}

func getAwsConfig(t *testing.T, ctx context.Context, endpoint string) aws.Config {
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

func createContainer(t *testing.T, ctx context.Context, containerRequest testcontainers.ContainerRequest) testcontainers.Container {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start LocalStack: %s", err)
	}
	return container
}

func createLocalstackContainerRequest() testcontainers.ContainerRequest {
	coreDynamoTablePath, _ := filepath.Abs("../../../../dev/localstack/cloudformation/core-dynamo-table.yml")
	initAwsPath, _ := filepath.Abs("../../../../dev/docker/scripts")
	return testcontainers.ContainerRequest{
		Image:        "localstack/localstack:1.3.0",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES":       "dynamodb,cloudformation",
			"DEFAULT_REGION": "us-east-1",
		},
		WaitingFor: wait.ForLog("Ready."),
		Mounts: []testcontainers.ContainerMount{
			testcontainers.BindMount(coreDynamoTablePath, "/opt/code/localstack/core-dynamo-table.yml"),
			testcontainers.BindMount(initAwsPath, "/etc/localstack/init/ready.d"),
			testcontainers.BindMount("/var/run/docker.sock", "/var/run/docker.sock"), // Añadir para Lambda
		},
	}
}

func clearDynamoTable(t *testing.T, dynamoDbClient *dynamodb.Client, tableName string) {
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
