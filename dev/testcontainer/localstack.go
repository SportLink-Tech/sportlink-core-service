package testcontainer

import (
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
)

func containerRequest() testcontainers.ContainerRequest {
	coreDynamoTablePath, _ := filepath.Abs("../../../../dev/localstack/cloudformation/core-dynamo-table.yml")
	sqsPath, _ := filepath.Abs("../../../../dev/localstack/cloudformation/sqs-queue.yml")
	initAwsPath, _ := filepath.Abs("../../../../dev/docker/scripts")
	return testcontainers.ContainerRequest{
		Image:        "localstack/localstack:1.3.0",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES":       "sqs,dynamodb,ssm,cf:cloudformation",
			"DEFAULT_REGION": "us-east-1",
			"AWS_REGION":     "us-east-1",
		},
		WaitingFor: wait.ForLog("Ready."),
		Mounts: []testcontainers.ContainerMount{
			testcontainers.BindMount(sqsPath, "/opt/code/localstack/sqs-queue.yml"),
			testcontainers.BindMount(coreDynamoTablePath, "/opt/code/localstack/core-dynamo-table.yml"),
			testcontainers.BindMount(initAwsPath, "/etc/localstack/init/ready.d"),
			testcontainers.BindMount("/var/run/docker.sock", "/var/run/docker.sock"),
		},
	}
}
